package main

import (
	"log"
	"net/http"
	"time"

	"github.com/envm-org/envm/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


type application struct {
	config config.Config
// logger 
// db drivers 
	
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("envm api is running on port " + app.config.Addr))
	})

	return r
}


func (app *application) run(h http.Handler) error {
	server := &http.Server{
		Addr: app.config.Addr,
		Handler: h,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout: 10 * time.Second,
	}
	
	log.Println("Starting server on port ", app.config.Addr)
	//TODO: implement graceful shutdown

	return server.ListenAndServe()
}
