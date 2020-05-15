package goback

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/robfig/cron/v3"
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
	//var secondParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional)
	//cron.NewParser()
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
	db, err := bolt.Open(filepath.Join(s.dbDir, s.appConfig.Name+".db"), 0755, &bolt.Options{Timeout: 1 * time.Second})
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
