package goback

import "github.com/devplayg/himma/v2"

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
	controller := NewController(s, &app)
	if err := controller.Start(); err != nil {
		log.Error(err)
	}
	return nil
}
