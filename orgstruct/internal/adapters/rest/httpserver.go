package rest

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// HttpConfig конфигурация для http сервера из переменных окружения.
type HttpConfig struct {
	Host         string `env:"HTTP_HOST"`
	Port         string `env:"HTTP_PORT"`
	ReadTimeout  int    `env:"HTTP_READ_TIMEOUT"`
	WriteTimeout int    `env:"HTTP_WRITE_TIMEOUT"`
	IdleTimeout  int    `env:"HTTP_IDLE_TIMEOUT"`
}

// Server хранит конфигурацию *http.Server.
type Server struct {
	srv *http.Server
}

// NewHttpServer конструктор для Server.
func NewHttpServer(handler http.Handler, c HttpConfig, log *slog.Logger) *Server {
	return &Server{
		srv: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", c.Host, c.Port),
			Handler:      handler,
			ReadTimeout:  time.Duration(c.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(c.WriteTimeout) * time.Second,
			IdleTimeout:  time.Duration(c.IdleTimeout) * time.Second,
			ErrorLog:     slog.NewLogLogger(log.Handler(), slog.LevelError),
		},
	}
}

// ListenAndServe запускает прослушивание запросов.
func (s *Server) ListenAndServe(c HttpConfig) error {
	slog.Info("start listening to HTTP", "host", c.Host, "port", c.Port)

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

// Shutdown останавливает прослушивание запросов.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
