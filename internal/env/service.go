package env

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	repo "github.com/envm-org/envm/internal/adapters/postgresql/sqlc"
)

type Service interface {
	CreateEnv(ctx context.Context, tempEnv repo.CreateEnvironmentParams) (repo.Environment, error)
	GetEnv(ctx context.Context, id pgtype.UUID) (repo.Environment, error)
	UpdateEnv(ctx context.Context, tempEnv repo.UpdateVariableParams) (repo.Variable, error)
	DeleteEnv(ctx context.Context, params repo.DeleteVariableParams) error
}

type svc struct {
	repo *repo.Queries
}

func NewService(repo *repo.Queries) Service {
	return &svc{repo: repo}
}

func (s *svc) CreateEnv(ctx context.Context, tempEnv repo.CreateEnvironmentParams) (repo.Environment, error) {
	return s.repo.CreateEnvironment(ctx, tempEnv)
}

func (s *svc) GetEnv(ctx context.Context, id pgtype.UUID) (repo.Environment, error) {
	return s.repo.GetEnvironment(ctx, id)
}

func (s *svc) UpdateEnv(ctx context.Context, tempEnv repo.UpdateVariableParams) (repo.Variable, error) {
	return s.repo.UpdateVariable(ctx, tempEnv)
}

func (s *svc) DeleteEnv(ctx context.Context, params repo.DeleteVariableParams) error {
	return s.repo.DeleteVariable(ctx, params)
}