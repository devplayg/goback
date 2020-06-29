package goback

func (s *Server) startHttpServer() error {
	controller := NewController(s, s.appConfig.HimmaConfig)
	if err := controller.Start(); err != nil {
		log.Error(err)
	}
	return nil
}
