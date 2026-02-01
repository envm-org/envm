package main

import (
	"log/slog"
	"os"

	"github.com/envm-org/envm/internal/config"
)


func main() {

	cfg := config.Config {
		DatabaseURI: "",
		Addr : ":8080",
	}

	api := application {
		config: cfg,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)


	h := api.mount()

	if err := api.run(h); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}

	
}
