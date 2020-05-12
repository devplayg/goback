package goback

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/devplayg/himma/v2"
	"github.com/devplayg/hippo/v2"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Server struct {
	hippo.Launcher // DO NOT REMOVE
	appConfig      *AppConfig
	config         *Config
	configFile     *os.File
	// dbFile         *os.File
	// tempDbFile     *os.File
	dbDir   string
	rwMutex *sync.RWMutex
	db      *bolt.DB
}

func NewServer(appConfig *AppConfig) *Server {
	return &Server{
		config: &Config{
			Storages: nil,
			Jobs:     nil,
		},
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

	//addr := s.config.Server.Address
	//if len(addr) < 1 {
	//	addr = ":8000"
	//}
	controller := NewController(s, &app)
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

	// if err := s.tempDbFile.Close(); err != nil {
	// 	return err
	// }
	// os.Remove(s.tempDbFile.Name())
	//
	// if err := s.dbFile.Close(); err != nil {
	// 	return err
	// }

	if err := s.db.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Server) init() error {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log = s.Log

	if err := s.initDirectories(); err != nil {
		return fmt.Errorf("failed to initialize directorie; %w", err)
	}

	if err := s.initDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database; %w", err)
	}

	if err := s.loadConfig(); err != nil {
		return fmt.Errorf("failed to load configuration; %w", err)
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
	db, err := bolt.Open(filepath.Join(s.dbDir, s.appConfig.Name+".db"), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	s.db = db
	return db.Batch(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(SummaryBucketName); err != nil {
			return fmt.Errorf("failed to create summary bucket; %w", err)
		}
		if _, err := tx.CreateBucketIfNotExists(BackupBucketName); err != nil {
			return fmt.Errorf("failed to create backup group bucket; %w", err)
		}
		if _, err := tx.CreateBucketIfNotExists(ConfigBucketName); err != nil {
			return fmt.Errorf("failed to create backup group bucket; %w", err)
		}
		return nil
	})
}

func (s *Server) issueDbId(bucketName []byte) (int, error) {
	var id int
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return ErrorBucketNotFound
		}
		newId, _ := b.NextSequence()
		id = int(newId)

		return b.Put(iToB(id), nil)
	})
	return id, err
}

func (s *Server) writeSummaries(results []*Summary) error {
	return s.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(SummaryBucketName)
		for i := range results {

			newSummaryId, _ := b.NextSequence()
			id := int(newSummaryId)
			results[i].Id = id
			data, err := results[i].Marshal()
			if err != nil {
				log.Error(err)
				continue
			}
			if err := b.Put(iToB(id), data); err != nil {
				log.Error(err)
				continue
			}
		}
		return nil
	})
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

	// Issue backup group id
	backupId, err := s.issueDbId(BackupBucketName)
	if err != nil {
		return fmt.Errorf("failed to issue backup id; %w", err)
	}

	go func() {
		started := time.Now()
		backup := NewBackup(backupId, job, s.dbDir, keeper, started)
		summaries, err := backup.Start()
		if err != nil {
			log.Error(err)
			return
		}

		if err := s.writeSummaries(summaries); err != nil {
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
	// get summary
	summary, err := s.findSummaryById(id)
	if err != nil {
		return nil, err
	}

	// read changes data in file
	h := md5.Sum([]byte(summary.SrcDir))
	key := hex.EncodeToString(h[:])
	logPath := filepath.Join(s.dbDir, fmt.Sprintf(ChangesDbName, key, summary.BackupId))
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return nil, err
	}
	return ioutil.ReadFile(logPath)
}
