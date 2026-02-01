package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	repo "github.com/envm-org/envm/internal/adapters/postgresql/sqlc"
)

type Service interface {
	CreateUser(ctx context.Context, params repo.CreateUserParams) (repo.User, error)
	GetUser(ctx context.Context, id pgtype.UUID) (repo.User, error)
	GetUserByEmail(ctx context.Context, email string) (repo.User, error)
	ListUsers(ctx context.Context) ([]repo.User, error)
	UpdateUser(ctx context.Context, params repo.UpdateUserParams) (repo.User, error)
	DeleteUser(ctx context.Context, id pgtype.UUID) error
}

type svc struct {
	repo *repo.Queries
}

func NewService(repo *repo.Queries) Service {
	return &svc{repo: repo}
}

func (s *svc) CreateUser(ctx context.Context, params repo.CreateUserParams) (repo.User, error) {
	return s.repo.CreateUser(ctx, params)
}

func (s *svc) GetUser(ctx context.Context, id pgtype.UUID) (repo.User, error) {
	return s.repo.GetUser(ctx, id)
}

func (s *svc) GetUserByEmail(ctx context.Context, email string) (repo.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *svc) ListUsers(ctx context.Context) ([]repo.User, error) {
	return s.repo.ListUsers(ctx)
}

func (s *svc) UpdateUser(ctx context.Context, params repo.UpdateUserParams) (repo.User, error) {
	return s.repo.UpdateUser(ctx, params)
}

func (s *svc) DeleteUser(ctx context.Context, id pgtype.UUID) error {
	return s.repo.DeleteUser(ctx, id)
}
