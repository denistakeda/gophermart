package server

import (
	"context"
	"gophermart/internal/core/services/logging"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Server struct {
	srv    *http.Server
	logger zerolog.Logger
}

func NewServer(addr string, handler http.Handler, log *logging.LoggerService) *Server {
	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		logger: log.ComponentLogger("Server"),
	}
}

func (s *Server) Start() {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal().Msgf("listen: %s\n", err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		s.logger.Fatal().Err(err).Msg("Server Shutdown")
	}
	s.logger.Info().Msg("Server exiting")
}
