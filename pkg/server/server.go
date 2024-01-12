package server

import (
	"context"
	"net/http"
	"time"

	"github.com/ermanimer/apigateway/pkg/config"
)

type Server struct {
	server          *http.Server
	shutdownTimeout time.Duration
}

func New(c config.Server) *Server {
	return &Server{
		server: &http.Server{
			Addr:           c.Address,
			ReadTimeout:    c.ReadTimeout,
			WriteTimeout:   c.WriteTimeout,
			IdleTimeout:    c.IdleTimeout,
			MaxHeaderBytes: c.MaxHeaderBytes,
			Handler:        http.NewServeMux(),
		},
		shutdownTimeout: c.ShutdownTimeout,
	}
}

func (s *Server) RegisterHandler(pattern string, handler http.Handler) {
	s.server.Handler.(*http.ServeMux).Handle(pattern, handler)
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
