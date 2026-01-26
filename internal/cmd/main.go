package main

import (
	"context"
	"log"
	"net/http"

	"github.com/envm-org/envm/internal/adapters/postgresql"
	repo "github.com/envm-org/envm/internal/adapters/postgresql/sqlc"
	"github.com/envm-org/envm/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	pool, err := postgresql.Connect(ctx, cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("Database connection failed: %v\n", err)
	}
	defer pool.Close()

	log.Println("Connected to database")

	queries := repo.New(pool)
	_ = queries // Keep compiler happy for now until we use it

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			http.Error(w, "Database not connected", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("OK"))
	})

	log.Println("Starting server on :5000")
	http.ListenAndServe(":5000", r)
}
