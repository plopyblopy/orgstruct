package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/plopyblopy/orgstruct/internal/adapters/postgres"
	"github.com/plopyblopy/orgstruct/internal/adapters/rest"
	"github.com/plopyblopy/orgstruct/internal/adapters/rest/handlers"
	"github.com/plopyblopy/orgstruct/internal/domain"
	"github.com/plopyblopy/orgstruct/internal/shared"
	"github.com/plopyblopy/orgstruct/internal/shared/config"
)

func main() {
	// config.
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	// logger.
	logLevelStr := cfg.LogLevel
	if logLevelStr == "" {
		logLevelStr = "info"
	}

	level, err := shared.ParseLevel(logLevelStr)
	if err != nil {
		level = slog.LevelInfo
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(log) // global logger.

	// service.
	err = app(*cfg, log)
	if err != nil {
		log.Error("app error", "err", err)
	}
	slog.Info("app stopped")
}

func app(cfg config.Config, log *slog.Logger) error {
	slog.Info("launching the application...")
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	// GORM
	// driver postgres.
	db := postgres.NewDb(cfg.DbConfig)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := recoverHandler(func() error { return db.Open(cfg.DbConfig) }); err != nil {
			errChan <- err
		}
	}()

	wg.Wait() // ожидает создания Db

	// router.
	router := rest.NewRouter()
	handlers.RegisterRoutes(router, db)
	handler := router.InitRouter(1)

	// http server.
	srv := rest.NewHttpServer(handler, cfg.HttpConfig, log)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// выполняется по srv.ListenAndServe
		// завершится только в виду err или Shutdown.
		if err := recoverHandler(func() error { return srv.ListenAndServe(cfg.HttpConfig) }); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// ожидает сиграл остановки приложения и реализует Graceful Shutdown.
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	slog.Info("the app is running!")

	select {
	case err := <-errChan:
		log.Error("server error", "err", err)
		return shutdownServer(srv, &wg, 5*time.Second)
	case <-stopChan:
		log.Info("shutdown signal received...")
		return shutdownServer(srv, &wg, 5*time.Second)
	}
}

// Graceful Shutdown для http server.
func shutdownServer(srv *rest.Server, wg *sync.WaitGroup, timeout time.Duration) error {
	slog.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("forced shutdown", "err", err)
		return err
	}

	wg.Wait()
	slog.Info("http server stopped")
	return nil
}

// recoverHandler функция-помощи для перехвата panic.
func recoverHandler(f func() error) error {
	err := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = domain.NewPanicError(fmt.Sprintf("%v", r))
			}
		}()

		return f()
	}()

	if err != nil {
		return err
	}

	return nil
}
