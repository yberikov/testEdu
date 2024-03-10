package server

import (
	"homework/internal/config"
	"homework/internal/handler"
	"homework/internal/middleware"
	"net/http"
)

func NewServer(cfg *config.Config, customMiddlewares ...func(http.Handler) http.Handler) *http.Server {
	mux := http.NewServeMux()

	deviceHandler := http.HandlerFunc(handler.GetDevice)

	mux.Handle("/device", deviceHandler)

	var muxWithMiddleware http.Handler
	muxWithMiddleware = mux

	for _, mw := range customMiddlewares {
		muxWithMiddleware = mw(muxWithMiddleware)
	}

	muxWithMiddleware = middleware.BasicAuthMiddleware(muxWithMiddleware)
	muxWithMiddleware = middleware.LoggingMiddleware(muxWithMiddleware)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      muxWithMiddleware,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	return srv
}
