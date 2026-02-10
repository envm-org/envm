package env

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	repo "github.com/envm-org/envm/internal/adapters/postgresql/sqlc"
	"github.com/envm-org/envm/pkg/encryption"
)

type Service interface {
	ListEnvs(ctx context.Context, projectID pgtype.UUID) ([]repo.Environment, error)
	CreateEnv(ctx context.Context, tempEnv repo.CreateEnvironmentParams) (repo.Environment, error)
	GetEnv(ctx context.Context, id pgtype.UUID) (repo.Environment, error)
	UpdateEnv(ctx context.Context, tempEnv repo.UpdateEnvironmentParams) (repo.Environment, error)
	DeleteEnv(ctx context.Context, id pgtype.UUID) error

	CreateVariable(ctx context.Context, params repo.CreateVariableParams) (repo.Variable, error)
	UpdateVariable(ctx context.Context, params repo.UpdateVariableParams) (repo.Variable, error)
	DeleteVariable(ctx context.Context, environmentID pgtype.UUID, key string) error
	ListVariables(ctx context.Context, environmentID pgtype.UUID) ([]repo.Variable, error)
}

type svc struct {
	repo          *repo.Queries
	encryptionKey string
}

func NewService(repo *repo.Queries, encryptionKey string) Service {
	return &svc{
		repo:          repo,
		encryptionKey: encryptionKey,
	}
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

func (s *svc) CreateVariable(ctx context.Context, params repo.CreateVariableParams) (repo.Variable, error) {
	encryptedValue, err := encryption.Encrypt(params.Value, s.encryptionKey)
	if err != nil {
		return repo.Variable{}, err
	}
	params.Value = encryptedValue
	return s.repo.CreateVariable(ctx, params)
}

func (s *svc) UpdateVariable(ctx context.Context, params repo.UpdateVariableParams) (repo.Variable, error) {
	encryptedValue, err := encryption.Encrypt(params.Value, s.encryptionKey)
	if err != nil {
		return repo.Variable{}, err
	}
	params.Value = encryptedValue
	return s.repo.UpdateVariable(ctx, params)
}

func (s *svc) DeleteVariable(ctx context.Context, environmentID pgtype.UUID, key string) error {
	return s.repo.DeleteVariable(ctx, repo.DeleteVariableParams{
		EnvironmentID: environmentID,
		Key:           key,
	})
}

func (s *svc) ListVariables(ctx context.Context, environmentID pgtype.UUID) ([]repo.Variable, error) {
	vars, err := s.repo.ListVariables(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	for i, v := range vars {
		decryptedValue, err := encryption.Decrypt(v.Value, s.encryptionKey)
		if err != nil {
			// If decryption fails, return the error or maybe the original value?
			// Returning error is safer to detect issues.
			return nil, err
		}
		vars[i].Value = decryptedValue
	}

	return vars, nil
}
