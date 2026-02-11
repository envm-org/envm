package project

import (
	"context"

	repo "github.com/envm-org/envm/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	ListProjects(ctx context.Context, userID pgtype.UUID) ([]repo.ListProjectsRow, error)
	CreateProject(ctx context.Context, tempProject repo.CreateProjectParams, userID pgtype.UUID) (repo.Project, error)
	GetProject(ctx context.Context, id pgtype.UUID) (repo.Project, error)
	UpdateProject(ctx context.Context, tempProject repo.UpdateProjectParams) (repo.Project, error)
	DeleteProject(ctx context.Context, id pgtype.UUID) error

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

func (s *svc) ListProjects(ctx context.Context, userID pgtype.UUID) ([]repo.ListProjectsRow, error) {
	return s.repo.ListProjects(ctx, userID)
}

func (s *svc) CreateProject(ctx context.Context, tempProject repo.CreateProjectParams, userID pgtype.UUID) (repo.Project, error) {
	project, err := s.repo.CreateProject(ctx, tempProject)
	if err != nil {
		return repo.Project{}, err
	}

	_, err = s.repo.AddProjectMember(ctx, repo.AddProjectMemberParams{
		ProjectID: project.ID,
		UserID:    userID,
		Role:      "owner",
	})
	if err != nil {
		return repo.Project{}, err
	}

	return project, nil
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
