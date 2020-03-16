package goback

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/devplayg/golibs/compress"
	"github.com/devplayg/golibs/converter"
	"github.com/devplayg/himma/v2"
	"github.com/devplayg/hippo/v2"
	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Server struct {
	hippo.Launcher // DO NOT REMOVE
	appConfig      *AppConfig
	config         *Config
	configFile     *os.File
	summaries      []*Summary
	dbFile         *os.File
	tempDbFile     *os.File
	dbDir          string
	rwMutex        *sync.RWMutex
}

func NewServer(appConfig *AppConfig) *Server {
	return &Server{
		appConfig: appConfig,
		dbDir:     "db",
		rwMutex:   new(sync.RWMutex),
	}
}

func NewEngine(appConfig *AppConfig) *hippo.Engine {
	server := NewServer(appConfig)
	engine := hippo.NewEngine(server, &hippo.Config{
		Name:        appConfig.Name,
		Description: appConfig.Description,
		Version:     appConfig.Version,
		Debug:       appConfig.Debug,
		Trace:       appConfig.Trace,
	})
	return engine
}

func (s *Server) Start() error {
	if err := s.init(); err != nil {
		return err
	}

	ch := make(chan struct{})
	go func() {
		if err := s.startHttpServer(); err != nil {
			s.Log.Error(err)
		}
		close(ch)
	}()

	defer func() {
		<-ch
	}()

	for {
		// Do your repetitive jobs
		// s.Log.Info("server is working on it")

		// Intentional error
		// s.Cancel() // send cancel signal to engine
		// return errors.New("intentional error")

		select {
		case <-s.Ctx.Done(): // for gracefully shutdown
			s.Log.Debug("server canceled; no longer works")
			return nil
		case <-time.After(2 * time.Second):
		}
	}
}

func (s *Server) startHttpServer() error {
	app := himma.Config{
		AppName:     s.appConfig.Name,
		Description: s.appConfig.Description,
		Url:         s.appConfig.Url,
		Phrase1:     s.appConfig.Text1,
		Phrase2:     s.appConfig.Text2,
		Year:        s.appConfig.Year,
		Version:     s.appConfig.Version,
		Company:     s.appConfig.Company,
	}

	addr := s.config.Server.Address
	if len(addr) < 1 {
		addr = ":8000"
	}
	controller := NewController(s, addr, &app)
	if err := controller.Start(); err != nil {
		log.Error(err)
	}
	return nil
}

func (s *Server) Stop() error {
	s.Log.Info("server has been stopped")
	if err := s.configFile.Close(); err != nil {
		return err
	}

	// if err := s.tempDbFile.Close(); err != nil {
	//	return err
	// }

	if err := s.tempDbFile.Close(); err != nil {
		return err
	}
	os.Remove(s.tempDbFile.Name())

	if err := s.dbFile.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Server) init() error {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log = s.Log

	if err := s.loadConfig(); err != nil {
		return fmt.Errorf("failed to load configuration; %w", err)
	}

	if err := s.initDirectories(); err != nil {
		return fmt.Errorf("failed to initialize directorie; %w", err)
	}

	if err := s.initDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database; %w", err)
	}

	return nil
}

func (s *Server) initDirectories() error {
	dbDir := filepath.Join(s.WorkingDir, s.dbDir)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.Mkdir(dbDir, 0600); err != nil {
			return fmt.Errorf("unable to create database directory: %w", err)
		}
	}
	s.dbDir = dbDir
	return nil
}

func (s *Server) initDatabase() error {

	// Open compressed database file
	path := filepath.Join(s.WorkingDir, "db", SummaryDbName)
	dbFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	s.dbFile = dbFile

	// Create temp database derived from compress database file
	tempDbPath := filepath.Join(s.WorkingDir, "db", SummaryTempDbName)
	tempDbFile, err := os.OpenFile(tempDbPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	s.tempDbFile = tempDbFile

	// Decompress compressed database file
	zr, err := gzip.NewReader(dbFile)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	data, err := ioutil.ReadAll(zr)
	if err != nil {
		return err
	}
	if _, err := tempDbFile.Write(data); err != nil {
		return err
	}
	if err := zr.Close(); err != nil {
		return err
	}

	return nil
}

// Thread-safe
func (s *Server) writeSummaries(results []*Summary) (int, error) {

	// Lock & unlock
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	s.tempDbFile.Seek(0, 0)
	data, err := ioutil.ReadAll(s.tempDbFile)
	if err != nil {
		return 0, err
	}

	// Decode gob-encoded data
	summaries, lastBackupId, lastSummaryId, err := DecodeSummaries(data)
	if err != nil {
		return 0, err
	}

	// Issue backup-id and summary-id
	backupId := lastBackupId + 1
	for i := range results {
		results[i].BackupId = backupId
		results[i].Id = lastSummaryId + 1
		lastSummaryId++
	}
	summaries = append(summaries, results...)

	// Encode data into gob-encoded data
	encoded, err := converter.EncodeToBytes(summaries)
	if err != nil {
		return 0, err
	}

	if err := s.tempDbFile.Truncate(0); err != nil {
		return 0, err
	}

	if _, err := s.tempDbFile.WriteAt(encoded, 0); err != nil {
		return 0, err
	}

	compressed, err := compress.Compress(encoded, compress.GZIP)
	if err != nil {
		return 0, fmt.Errorf("failed to compress summary data: %w", err)
	}

	if err := s.dbFile.Truncate(0); err != nil {
		return 0, err
	}

	if _, err := s.dbFile.WriteAt(compressed, 0); err != nil {
		return 0, err
	}

	// Gob decode

	// Issue backup-id and summary-id
	// Append summaries
	// Save
	return backupId, nil
}

// Thread-safe
func (s *Server) getSummaries() ([]*Summary, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	s.tempDbFile.Seek(0, 0)
	data, err := ioutil.ReadAll(s.tempDbFile)
	if err != nil {
		return nil, err
	}

	summaries, _, _, err := DecodeSummaries(data)
	if err != nil {
		return nil, err
	}

	return summaries, nil
}

//
// func (s *Server) findSummaries() ([]*Summary, error){
//	var summaries []*Summary
//	json.Unmarshal(s.tempDbFile, &summaries)
//	return nil, nil
// }

func (s *Server) loadConfig() error {
	file, err := os.OpenFile(ConfigFileName, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	s.configFile = file

	rows := make([]string, 0)
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		rows = append(rows, fileScanner.Text())
	}

	if err := yaml.Unmarshal([]byte(strings.Join(rows, "\n")), &s.config); err != nil {
		return err
	}
	return nil
}

func (s *Server) saveConfig() error {
	if err := s.configFile.Truncate(0); err != nil {
		return err
	}

	data, err := yaml.Marshal(s.config)
	if err != nil {
		return err
	}

	if _, err := s.configFile.WriteAt(data, 0); err != nil {
		return err
	}

	return nil
}

func (s *Server) runBackupJob(jobId int) error {
	job := s.config.findJobById(jobId)
	if job == nil {
		return errors.New("backup job not found")
	}

	keeper := NewKeeper(job)
	if keeper == nil {
		return fmt.Errorf("invalid keeper protocol %d", job.Storage.Protocol)
	}

	go func() {
		started := time.Now()
		backup := NewBackup(job, s.dbDir, keeper, started, s.appConfig.Debug)
		summaries, err := backup.Start()
		if err != nil {
			log.Error(err)
			return
		}

		backupId, err := s.writeSummaries(summaries)
		if err != nil {
			log.Error(err)
			return
		}

		log.WithFields(logrus.Fields{
			"execTime": time.Since(started).Seconds(),
			"backupId": backupId,
		}).Info("## all backup processes done")
	}()

	return nil
}

func (s *Server) getChangesLog(id int) ([]byte, error) { // wondory
	//	summary := c.findSummaryById(id)
	//	if summary == nil {
	//		return nil, errors.New("summary not found")
	//	}
	//
	//	h := md5.Sum([]byte(summary.SrcDir))
	//	key := hex.EncodeToString(h[:])
	//	logPath := filepath.Join(c.dbDir, fmt.Sprintf(ChangesDbName, key, summary.BackupId))
	//	if _, err := os.Stat(logPath); os.IsNotExist(err) {
	//		return nil, err
	//	}
	//
	//	return ioutil.ReadFile(logPath)
	return nil, nil
}
