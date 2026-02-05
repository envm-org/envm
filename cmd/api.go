package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	repo "github.com/envm-org/envm/internal/adapters/postgresql/sqlc"
	"github.com/envm-org/envm/internal/auth"
	"github.com/envm-org/envm/internal/env"
	appMiddleware "github.com/envm-org/envm/internal/middleware"
	"github.com/envm-org/envm/internal/org"
	"github.com/envm-org/envm/internal/project"
	"github.com/envm-org/envm/internal/users"
	authPkg "github.com/envm-org/envm/pkg/auth"
	"github.com/envm-org/envm/pkg/email"
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

	emailSender := email.NewLogSender()

	// Authorizer
	authorizer := auth.NewAuthorizer(q)

	// Env
	envService := env.NewService(q)
	envHandler := env.NewHandler(envService, authorizer)

	// Project
	projectService := project.NewService(q)
	projectHandler := project.NewHandler(projectService, authorizer)

	// Org
	orgService := org.NewService(q, emailSender)
	orgHandler := org.NewHandler(orgService, authorizer)

	// Auth
	tokenMaker, err := authPkg.NewJWTMaker(app.config.TokenSecret)
	if err != nil {
		panic(fmt.Errorf("cannot create token maker: %w", err))
	}
	authMiddleware := appMiddleware.AuthMiddleware(tokenMaker)

	authService := auth.NewService(q, emailSender)
	authHandler := auth.NewHandler(authService, tokenMaker)

	// Users
	usersService := users.NewService(q)
	usersHandler := users.NewHandler(usersService)

	// Public Routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.Register)
	})

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Route("/users", func(r chi.Router) {
			r.Get("/", usersHandler.GetUser)
			r.Put("/", usersHandler.UpdateUser)
			r.Delete("/", usersHandler.DeleteUser)
			r.Get("/list", usersHandler.ListUsers)
			r.Get("/email", usersHandler.GetUserByEmail)
		})

		r.Route("/env", func(r chi.Router) {
			r.Post("/", envHandler.CreateEnv)
			r.Get("/", envHandler.GetEnv)
			r.Put("/", envHandler.UpdateEnv)
			r.Delete("/", envHandler.DeleteEnv)
			r.Get("/list", envHandler.ListEnvs)
		})

		r.Route("/project", func(r chi.Router) {
			r.Post("/", projectHandler.CreateProject)
			r.Get("/", projectHandler.GetProject)
			r.Put("/", projectHandler.UpdateProject)
			r.Delete("/", projectHandler.DeleteProject)
			r.Post("/members", projectHandler.AddMember)
			r.Delete("/members", projectHandler.RemoveMember)
			r.Get("/members", projectHandler.ListMembers)
		})

		r.Route("/org", func(r chi.Router) {
			r.Post("/", orgHandler.CreateOrg)
			r.Get("/", orgHandler.GetOrg)
			r.Put("/", orgHandler.UpdateOrg)
			r.Delete("/", orgHandler.DeleteOrg)
			r.Get("/list", orgHandler.ListOrgs)

			r.Post("/invite", orgHandler.InviteMember)
			r.Post("/join", orgHandler.AcceptInvitation)
		})
	})

	return r
}

func (app *application) run(h http.Handler) error {
	server := &http.Server{
		Addr:         app.config.Addr,
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	log.Println("Starting server on port ", app.config.Addr)
	//TODO: implement graceful shutdown

	return server.ListenAndServe()
}
