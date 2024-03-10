package server

import (
	"github.com/yberikov/configLoader"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg *configLoader.Config, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	}

	return s.httpServer.ListenAndServe()
}
