package goback

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/devplayg/hippo/v2"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func NewEngine(app *AppConfig) *hippo.Engine {
	server := NewServer(app)
	logDir := "."
	if app.Verbose {
		logDir = ""
	}
	engine := hippo.NewEngine(server, &hippo.Config{
		Name:        app.Name,
		Description: app.Description,
		Version:     app.Version,
		Debug:       app.Debug,
		Trace:       app.Trace,
		LogDir:      logDir,
	})
	return engine
}

func NewServer(appConfig *AppConfig) *Server {
	return &Server{
		dbDir:     "db",
		appConfig: appConfig,
		config:    NewConfig(),
		rwMutex:   new(sync.RWMutex),
	}
}

type Server struct {
	hippo.Launcher // DO NOT REMOVE
	appConfig      *AppConfig
	config         *Config
	configFile     *os.File
	dbDir          string
	rwMutex        *sync.RWMutex
	db             *bolt.DB
	cron           *cron.Cron
}

func (s *Server) Start() error {
	// Initialize server
	if err := s.init(); err != nil {
		return err
	}

	ch := make(chan struct{})

	// Start HTTP server
	go func() {
		defer close(ch)
		if err := s.startHttpServer(); err != nil {
			s.Log.Error(err)
			return
		}
	}()

	// Start scheduler
	s.cron.Start()

	log.Infof("server has been started; listening on %s", s.appConfig.Address)

	// Wait for HTTP server to stop
	<-ch

	return nil
}

func (s *Server) Stop() error {

	if s.cron != nil {
		s.cron.Stop()
	}

	if s.db != nil {
		if err := s.db.Close(); err != nil {
			return err
		}
	}

	s.Log.Info("server has been stopped")
	return nil
}

func (s *Server) runBackupJob(jobId int) error {

	// Get job
	job := s.findJobById(jobId)
	if job == nil {
		return fmt.Errorf("job-%d not found", jobId)
	}
	if job.running {
		return fmt.Errorf("job-%d is already running now", job.Id)
	}

	// Get storage
	job.Storage = s.findStorageById(job.StorageId)
	if job.Storage == nil {
		return fmt.Errorf("storage-%d not found", job.StorageId)
	}

	// Get keeper
	keeper := NewKeeper(job)
	if keeper == nil {
		return fmt.Errorf("invalid keeper protocol %d", job.Storage.Protocol)
	}

	// log.WithFields(logrus.Fields{
	//	"id":      job.Id,
	//	"running": job.running,
	// }).Debug("job")

	s.rwMutex.Lock()
	job.running = true
	s.rwMutex.Unlock()

	go func() {
		defer func() {
			job.running = false
		}()

		// Issue backup group id
		backupId, err := s.issueDbId(BackupBucket)
		if err != nil {
			log.Error(fmt.Errorf("failed to issue backup id; %w", err))
			return
		}

		started := time.Now()
		backup := NewBackup(backupId, job, s.dbDir, keeper, started)
		summaries, err := backup.Start()
		if err != nil {
			log.WithFields(logrus.Fields{
				"jobId":      job.Id,
				"backupType": job.BackupType,
				"storageId":  job.StorageId,
			}).Error(err)
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

func (s *Server) getChangesLog(id int) ([]byte, error) {
	// get summary
	summary, err := s.findSummaryById(id)
	if err != nil {
		return nil, err
	}

	// read changes data in file
	h := md5.Sum([]byte(summary.SrcDir))
	key := hex.EncodeToString(h[:])
	logPath := filepath.Join(s.dbDir, fmt.Sprintf(ChangesDbName, summary.BackupId, key))
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("database '%s' not found", filepath.Base(logPath))
	}
	return ioutil.ReadFile(logPath)
}
