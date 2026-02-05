package org

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	repo "github.com/envm-org/envm/internal/adapters/postgresql/sqlc"
	"github.com/envm-org/envm/pkg/email"
)

type Service interface {
	CreateOrg(ctx context.Context, params repo.CreateOrganizationParams) (repo.Organization, error)
	GetOrg(ctx context.Context, id pgtype.UUID) (repo.Organization, error)
	ListOrgs(ctx context.Context) ([]repo.Organization, error)
	UpdateOrg(ctx context.Context, params repo.UpdateOrganizationParams) (repo.Organization, error)
	DeleteOrg(ctx context.Context, id pgtype.UUID) error
	
	InviteMember(ctx context.Context, orgID pgtype.UUID, email, role string, invitedBy pgtype.UUID) error
	AcceptInvitation(ctx context.Context, token string, userID pgtype.UUID) error
}

type svc struct {
	repo   *repo.Queries
	mailer email.Sender
}

func NewService(repo *repo.Queries, mailer email.Sender) Service {
	return &svc{
		repo:   repo,
		mailer: mailer,
	}
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

func (s *svc) UpdateOrg(ctx context.Context, params repo.UpdateOrganizationParams) (repo.Organization, error) {
	return s.repo.UpdateOrganization(ctx, params)
}

func (s *svc) DeleteOrg(ctx context.Context, id pgtype.UUID) error {
	return s.repo.DeleteOrganization(ctx, id)
}

func (s *svc) InviteMember(ctx context.Context, orgID pgtype.UUID, emailAddr, role string, invitedBy pgtype.UUID) error {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return err
	}
	token := hex.EncodeToString(tokenBytes)
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days

	_, err := s.repo.CreateInvitation(ctx, repo.CreateInvitationParams{
		OrganizationID: orgID,
		Email:          emailAddr,
		Role:           role,
		Token:          token,
		ExpiresAt:      pgtype.Timestamptz{Time: expiresAt, Valid: true},
		InvitedBy:      invitedBy,
	})
	if err != nil {
		return fmt.Errorf("failed to create invitation: %w", err)
	}

	// Send Email
	subject := "You are invited to join an organization"
	body := fmt.Sprintf("You have been invited to join. Use this token to accept: %s", token)
	
	return s.mailer.SendEmail(emailAddr, subject, body)
}

func (s *svc) AcceptInvitation(ctx context.Context, token string, userID pgtype.UUID) error {
	invitation, err := s.repo.GetInvitationByToken(ctx, token)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("invalid or expired invitation")
		}
		return err
	}

	_, err = s.repo.AddOrganizationMember(ctx, repo.AddOrganizationMemberParams{
		OrganizationID: invitation.OrganizationID,
		UserID:         userID,
		Role:           invitation.Role,
	})
	if err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	err = s.repo.DeleteInvitation(ctx, invitation.ID)
	if err != nil {
		return err
	}

	return nil
}
