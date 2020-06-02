package goback

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

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

	if err := s.initScheduler(); err != nil {
		return fmt.Errorf("failed to initialize scheduler; %w", err)
	}
	return nil
}

func (s *Server) initScheduler() error {
	loc := time.Local
	s.cron = cron.New(cron.WithLocation(loc))

	for i, job := range s.config.Jobs {
		if !job.Enabled {
			continue
		}
		if job.SrcDirs == nil || len(job.SrcDirs) < 1 {
			continue
		}

		if len(job.Schedule) < 1 {
			continue
		}

		jobId := s.config.Jobs[i].Id
		entryId, err := s.cron.AddFunc(s.config.Jobs[i].Schedule, func() {
			log.WithFields(logrus.Fields{
				"jobId": jobId,
			}).Info("RUN SCHEDULER")
			if err := s.runBackupJob(jobId); err != nil {
				log.Error(err)
			}
		})
		if err != nil {
			return fmt.Errorf("failed to load scheduler %d; %w", job.Id, err)
		}
		s.config.Jobs[i].cronEntryId = &entryId
		log.WithFields(logrus.Fields{
			"jobId":     job.Id,
			"scheduler": job.Schedule,
		}).Info("backup scheduler loaded")

	}

	//var secondParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional)
	//cron.NewParser()
	return nil
}

func (s *Server) initDirectories() error {
	dbDir := filepath.Join(s.WorkingDir, s.dbDir)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.Mkdir(dbDir, 0755); err != nil {
			return fmt.Errorf("unable to create database directory: %w", err)
		}
	}
	s.dbDir = dbDir
	return nil
}

func (s *Server) initDatabase() error {
	db, err := bolt.Open(filepath.Join(s.dbDir, s.appConfig.Name+".db"), 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	s.db = db
	return db.Batch(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(SummaryBucket); err != nil {
			return fmt.Errorf("failed to create summary bucket; %w", err)
		}
		if _, err := tx.CreateBucketIfNotExists(BackupBucket); err != nil {
			return fmt.Errorf("failed to create backup group bucket; %w", err)
		}
		if _, err := tx.CreateBucketIfNotExists(ConfigBucket); err != nil {
			return fmt.Errorf("failed to create backup group bucket; %w", err)
		}
		return nil
	})
}
