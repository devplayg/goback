package goback

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/devplayg/himma"
	"github.com/devplayg/hippo/v2"
	log "github.com/sirupsen/logrus"
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

	spew.Dump(s.config.Jobs)

	// if err := s.startHttpServer(); err != nil {
	// 	return err
	// }
	//
	// for {
	// 	s.Log.Info("server is working on it")
	//
	// 	// return errors.New("intentional error")
	//
	// 	select {
	// 	case <-s.Ctx.Done(): // for gracefully shutdown
	// 		return nil
	// 	case <-time.After(2 * time.Second):
	// 	}
	// }

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
	app := himma.Application{
		AppName:     s.appConfig.Name,
		Description: s.appConfig.Description,
		Url:         s.appConfig.Url,
		Phrase1:     s.appConfig.Text1,
		Phrase2:     s.appConfig.Text2,
		Year:        s.appConfig.Year,
		Version:     s.appConfig.Version,
		Company:     s.appConfig.Company,
	}
	controller := NewController(s, "db", s.appConfig.Addr, &app)
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

	return nil
}
