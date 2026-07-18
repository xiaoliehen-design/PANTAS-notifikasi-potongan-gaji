package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bcpriok/pantas/internal/auth"
	"github.com/bcpriok/pantas/internal/config"
	"github.com/bcpriok/pantas/internal/database"
	"github.com/bcpriok/pantas/internal/httpapi"
	"github.com/bcpriok/pantas/internal/importer"
	"github.com/bcpriok/pantas/internal/mailer"
	"github.com/bcpriok/pantas/internal/storage"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
	cfg, err := config.Load()
	if err != nil {
		logger.Error("configuration invalid", "error", err)
		os.Exit(1)
	}

	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	pool, err := database.Open(rootCtx, cfg)
	if err != nil {
		logger.Error("database unavailable", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	if err := database.BootstrapAdmin(rootCtx, pool, cfg.BootstrapAdminUsername, cfg.BootstrapAdminName, cfg.BootstrapAdminPassword); err != nil {
		logger.Error("bootstrap admin failed; pastikan migration sudah dijalankan", "error", err)
		os.Exit(1)
	}

	authService := auth.New(pool, cfg)
	importService := importer.New(pool)
	storageClient := storage.New(cfg)
	app, err := httpapi.New(pool, cfg, authService, importService, storageClient, logger)
	if err != nil {
		logger.Error("build web app", "error", err)
		os.Exit(1)
	}
	worker := mailer.New(pool, cfg, logger)
	go worker.Run(rootCtx)

	server := &http.Server{
		Addr: cfg.Address(), Handler: app.Handler(),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       3 * time.Minute,
		WriteTimeout:      3 * time.Minute,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	go func() {
		logger.Info("PANTAS started", "config", cfg.String())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server stopped", "error", err)
			stop()
		}
	}()

	<-rootCtx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown", "error", err)
	}
}
