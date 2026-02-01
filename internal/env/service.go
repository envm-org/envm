package env

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	repo "github.com/envm-org/envm/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListEnvs(ctx context.Context, projectID pgtype.UUID) ([]repo.Environment, error)
	CreateEnv(ctx context.Context, tempEnv repo.CreateEnvironmentParams) (repo.Environment, error)
	GetEnv(ctx context.Context, id pgtype.UUID) (repo.Environment, error)
	UpdateEnv(ctx context.Context, tempEnv repo.UpdateEnvironmentParams) (repo.Environment, error)
	DeleteEnv(ctx context.Context, id pgtype.UUID) error
}

type svc struct {
	repo *repo.Queries
}

func NewService(repo *repo.Queries) Service {
	return &svc{repo: repo}
}

func (s *svc) ListEnvs(ctx context.Context, projectID pgtype.UUID) ([]repo.Environment, error) {
	return s.repo.ListEnvironments(ctx, projectID)
}

func (s *svc) CreateEnv(ctx context.Context, tempEnv repo.CreateEnvironmentParams) (repo.Environment, error) {
	return s.repo.CreateEnvironment(ctx, tempEnv)
}

func (s *svc) GetEnv(ctx context.Context, id pgtype.UUID) (repo.Environment, error) {
	return s.repo.GetEnvironment(ctx, id)
}

func (s *svc) UpdateEnv(ctx context.Context, tempEnv repo.UpdateEnvironmentParams) (repo.Environment, error) {
	return s.repo.UpdateEnvironment(ctx, tempEnv)
}

func (s *svc) DeleteEnv(ctx context.Context, id pgtype.UUID) error {
	return s.repo.DeleteEnvironment(ctx, id)
}