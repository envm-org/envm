package main

import (
	"fmt"
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
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
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

	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(appMiddleware.SecurityHeaders)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("envm api is running on port " + app.config.Addr))
	})

	q := repo.New(app.db)

	emailSender := email.NewLogSender()

	authorizer := auth.NewAuthorizer(q)

	envService := env.NewService(q)
	envHandler := env.NewHandler(envService, authorizer)

	projectService := project.NewService(q)
	projectHandler := project.NewHandler(projectService, authorizer)

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

	usersService := users.NewService(q)
	usersHandler := users.NewHandler(usersService)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.Register)
		r.Post("/refresh", authHandler.Refresh)
		r.Post("/logout", authHandler.Logout)
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
