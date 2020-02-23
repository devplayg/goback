package goback

import (
	"github.com/devplayg/himma/v2"
	"github.com/devplayg/hippo/v2"
	"runtime"
	"time"
)

type Server struct {
	hippo.Launcher // DO NOT REMOVE
	appConfig      *AppConfig
	config         *Config
}

func NewServer(appConfig *AppConfig) *Server {
	return &Server{
		appConfig: appConfig,
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

	for _, job := range s.config.Jobs {

		// Local backup
		if job.Storage.Protocol == LocalDisk {
			// log.WithFields(logrus.Fields{
			// 	"target": "localDisk",
			// 	"dir":    job.Storage.Dir,
			// }).Debug("backup")
			backup := NewBackup(job.SrcDirs, NewLocalKeeper(job.Storage.Dir), job.BackupType, s.appConfig.Debug)
			if err := backup.Start(); err != nil {
				log.Error(err)
			}
			continue
		}

		if job.Storage.Protocol == Sftp {
			backup := NewBackup(job.SrcDirs, NewSftpKeeper(job.Storage), job.BackupType, s.appConfig.Debug)
			if err := backup.Start(); err != nil {
				log.Error(err)
			}
			continue
		}
	}

	for {
		// Do your repetitive jobs
		s.Log.Info("server is working on it")

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
	controller := NewController(s, "db", addr, &app)
	if err := controller.Start(); err != nil {
		log.Error(err)
	}
	return nil
}

func (s *Server) Stop() error {
	s.Log.Info("server has been stopped")
	return nil
}

func (s *Server) init() error {
	config, err := loadConfig(ConfigFileName)
	if err != nil {
		return err
	}
	s.config = config

	runtime.GOMAXPROCS(runtime.NumCPU())
	log = s.Log

	return nil
}
