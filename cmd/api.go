package main

import (
	"log"
	"net/http"
	"time"

	repo "github.com/envm-org/envm/internal/adapters/postgresql/sqlc"
	"github.com/envm-org/envm/internal/env"
	"github.com/envm-org/envm/internal/org"
	"github.com/envm-org/envm/internal/project"
	"github.com/envm-org/envm/internal/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)


type application struct {
	config Config
	db     *pgx.Conn
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

	q := repo.New(app.db)

	// Env
	envService := env.NewService(q)
	envHandler := env.NewHandler(envService)
	r.Route("/env", func(r chi.Router) {
		r.Post("/", envHandler.CreateEnv)
		r.Get("/", envHandler.GetEnv)
		r.Put("/", envHandler.UpdateEnv)
		r.Delete("/", envHandler.DeleteEnv)
		r.Get("/list", envHandler.ListEnvs)
	})

	// Project
	projectService := project.NewService(q)
	projectHandler := project.NewHandler(projectService)
	r.Route("/project", func(r chi.Router) {
		r.Post("/", projectHandler.CreateProject)
		r.Get("/", projectHandler.GetProject)
		r.Put("/", projectHandler.UpdateProject)
		r.Delete("/", projectHandler.DeleteProject)
		r.Get("/list", projectHandler.ListProjects)
	})

	// Org
	orgService := org.NewService(q)
	orgHandler := org.NewHandler(orgService)
	r.Route("/org", func(r chi.Router) {
		r.Post("/", orgHandler.CreateOrg)
		r.Get("/", orgHandler.GetOrg)
		r.Put("/", orgHandler.UpdateOrg)
		r.Delete("/", orgHandler.DeleteOrg)
		r.Get("/list", orgHandler.ListOrgs)
	})

	// Users
	usersService := users.NewService(q)
	usersHandler := users.NewHandler(usersService)
	r.Route("/users", func(r chi.Router) {
		r.Post("/", usersHandler.CreateUser)
		r.Get("/", usersHandler.GetUser)
		r.Put("/", usersHandler.UpdateUser)
		r.Delete("/", usersHandler.DeleteUser)
		r.Get("/list", usersHandler.ListUsers)
		r.Get("/email", usersHandler.GetUserByEmail)
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
