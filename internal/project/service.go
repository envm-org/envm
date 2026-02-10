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

	ListProjectsForMember(ctx context.Context, organizationID, userID pgtype.UUID) ([]repo.Project, error)
	AddMember(ctx context.Context, projectID, userID pgtype.UUID, role string) error
	RemoveMember(ctx context.Context, projectID, userID pgtype.UUID) error
	GetMember(ctx context.Context, projectID, userID pgtype.UUID) (repo.ProjectMember, error)
	ListMembers(ctx context.Context, projectID pgtype.UUID) ([]repo.ListProjectMembersRow, error)
}

type svc struct {
	repo *repo.Queries
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

func (s *svc) ListProjectsForMember(ctx context.Context, organizationID, userID pgtype.UUID) ([]repo.Project, error) {
	projects, err := s.repo.ListProjectsForMember(ctx, repo.ListProjectsForMemberParams{
		OrganizationID: organizationID,
		UserID:         userID,
	})
	if err != nil {
		return nil, err
	}

	result := make([]repo.Project, len(projects))
	for i, p := range projects {
		result[i] = repo.Project{
			ID:             p.ID,
			OrganizationID: p.OrganizationID,
			Name:           p.Name,
			Slug:           p.Slug,
			Description:    p.Description,
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
		}
	}
	return result, nil
}

func (s *svc) AddMember(ctx context.Context, projectID, userID pgtype.UUID, role string) error {
	_, err := s.repo.AddProjectMember(ctx, repo.AddProjectMemberParams{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
	})
	return err
}

func (s *svc) RemoveMember(ctx context.Context, projectID, userID pgtype.UUID) error {
	return s.repo.RemoveProjectMember(ctx, repo.RemoveProjectMemberParams{
		ProjectID: projectID,
		UserID:    userID,
	})
}

func (s *svc) GetMember(ctx context.Context, projectID, userID pgtype.UUID) (repo.ProjectMember, error) {
	return s.repo.GetProjectMember(ctx, repo.GetProjectMemberParams{
		ProjectID: projectID,
		UserID:    userID,
	})
}

func (s *svc) ListMembers(ctx context.Context, projectID pgtype.UUID) ([]repo.ListProjectMembersRow, error) {
	return s.repo.ListProjectMembers(ctx, projectID)
}
