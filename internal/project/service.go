package project

import (
	"context"

	repo "github.com/envm-org/envm/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	ListProjects(ctx context.Context, organizationID pgtype.UUID) ([]repo.Project, error)
	CreateProject(ctx context.Context, tempProject repo.CreateProjectParams) (repo.Project, error)
	GetProject(ctx context.Context, id pgtype.UUID) (repo.Project, error)
	UpdateProject(ctx context.Context, tempProject repo.UpdateProjectParams) (repo.Project, error)
	DeleteProject(ctx context.Context, id pgtype.UUID) error
}

type svc struct {
	repo  *repo.Queries
}

func NewService(repo *repo.Queries) Service {
	return &svc{repo: repo}
}


func (s *svc) ListProjects(ctx context.Context, organizationID pgtype.UUID) ([]repo.Project, error) {
	return s.repo.ListProjects(ctx, organizationID)
}


func (s *svc) CreateProject(ctx context.Context, tempProject repo.CreateProjectParams) (repo.Project, error) {
	return s.repo.CreateProject(ctx, tempProject)
}


func (s *svc) GetProject(ctx context.Context, id pgtype.UUID) (repo.Project, error) {
	return s.repo.GetProject(ctx, id)
}


func (s *svc) UpdateProject(ctx context.Context, tempProject repo.UpdateProjectParams) (repo.Project, error) {
	return s.repo.UpdateProject(ctx, tempProject)
}


func (s *svc) DeleteProject(ctx context.Context, id pgtype.UUID) error {
	return s.repo.DeleteProject(ctx, id)
}
