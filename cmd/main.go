package main

import (
	"context"
	"log/slog"
	"os"

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
	
	api := application {
		config: cfg,
		db: conn,
	}

	h := api.mount()

	if err := api.run(h); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}

	
}
