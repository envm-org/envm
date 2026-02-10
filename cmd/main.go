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

	"github.com/envm-org/envm/pkg/env"
	"github.com/jackc/pgx/v5"
)

type Config struct {
	DatabaseURI string
	Addr        string
	TokenSecret string
	DB          DBConfig
}

type DBConfig struct {
	DSN string
}

func main() {
	ctx := context.Background()

	addr := env.GetString("ADDR", ":8080")
	cfg := Config{
		DatabaseURI: env.GetString("DATABASE_URI", "postgres://postgres:postgres@localhost:5432/envm"),
		Addr:        addr,
		TokenSecret: env.GetString("TOKEN_SECRET", "12345678901234567890123456789012"),
	}
	cfg.DB.DSN = cfg.DatabaseURI

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// database connection
	conn, err := pgx.Connect(ctx, cfg.DB.DSN)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		panic(err)
	}
	defer conn.Close(ctx)

	logger.Info("database connection successful")

	api := application{
		config: cfg,
		db:     conn,
	}

	h := api.mount()

	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("starting server", "addr", cfg.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logger.Error("server error", "error", err)
		os.Exit(1)

	case sig := <-shutdown:
		logger.Info("shutdown signal received", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("graceful shutdown failed", "error", err)
			// Force close if graceful shutdown fails
			if err := server.Close(); err != nil {
				logger.Error("failed to close server", "error", err)
			}
			os.Exit(1)
		}

		// Check if context deadline was exceeded
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			logger.Warn("shutdown deadline exceeded, some requests may have been terminated")
		}

		logger.Info("server stopped gracefully")
	}
}
