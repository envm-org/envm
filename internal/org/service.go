package org

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	repo "github.com/envm-org/envm/internal/adapters/postgresql/sqlc"
)

type Service interface {
	CreateOrg(ctx context.Context, params repo.CreateOrganizationParams) (repo.Organization, error)
	GetOrg(ctx context.Context, id pgtype.UUID) (repo.Organization, error)
	ListOrgs(ctx context.Context) ([]repo.Organization, error)
}

type svc struct {
	repo *repo.Queries
}

func NewService(repo *repo.Queries) Service {
	return &svc{repo: repo}
}

func (s *svc) CreateOrg(ctx context.Context, params repo.CreateOrganizationParams) (repo.Organization, error) {
	return s.repo.CreateOrganization(ctx, params)
}

func (s *svc) GetOrg(ctx context.Context, id pgtype.UUID) (repo.Organization, error) {
	return s.repo.GetOrganization(ctx, id)
}

func (s *svc) ListOrgs(ctx context.Context) ([]repo.Organization, error) {
	return s.repo.ListOrganizations(ctx)
}
