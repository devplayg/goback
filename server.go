package goback

import (
	"github.com/devplayg/himma"
	"github.com/devplayg/hippo/v2"
	"time"
)

type Server struct {
	hippo.Launcher // DO NOT REMOVE
	himmaConfig    *Config
	addr           string
}

func NewServer(himmaConfig himma.Application addr string) *Server {
	return &Server{
	}
}

func (s *Server) Start() error {
	if err := s.init(); err != nil {
		return err
	}

	if err != s.startWeb(); err != nil {
		return err
	}

	for {
		s.Log.Info("server is working on it")

		// return errors.New("intentional error")

		select {
		case <-s.Done: // for gracefully shutdown
			return nil
		case <-time.After(2 * time.Second):
		}
	}
}

func (s *Server) startWeb() error {
	//app := himma.Application{
	//	AppName:     "SecuBACKUP",
	//	Description: "INCREMENTAL BACKUP ",
	//	Url:         "https://devplayg.com",
	//	Phrase1:     "KEEP YOUR DATA SAFE",
	//	Phrase2:     "Powered by Go",
	//	Year:        time.Now().Year(),
	//	Version:     appVersion,
	//	Company:     "SECUSOLUTION",
	//}
	//c := goback.NewController(backup.DbDir, "127.0.0.1:8000", &app)
	//if err := c.Start(); err != nil {
	//	log.Error(err)
	//}
}

func (s *Server) Stop() error {
	s.Log.Info("server has been stopped")
	return nil
}

func (s *Server) init() error {
	//config, err := loadConfig(s.)
	//if err != nil {
	//	return err
	//}
	//s.config = config
	return nil
}
