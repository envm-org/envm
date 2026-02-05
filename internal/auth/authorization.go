package auth

import (
	"context"
	"fmt"
	"slices"

	repo "github.com/envm-org/envm/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Role string

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

type Authorizer interface {
	HasRole(ctx context.Context, userID, orgID pgtype.UUID, requiredRoles ...Role) error
}

type authorizer struct {
	repo *repo.Queries
}

func NewAuthorizer(repo *repo.Queries) Authorizer {
	return &authorizer{repo: repo}
}

func (a *authorizer) HasRole(ctx context.Context, userID, orgID pgtype.UUID, requiredRoles ...Role) error {
	member, err := a.repo.GetOrganizationMember(ctx, repo.GetOrganizationMemberParams{
		OrganizationID: orgID,
		UserID:         userID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user is not a member of this organization")
		}
		return fmt.Errorf("failed to check membership: %w", err)
	}

	for _, role := range requiredRoles {
		if Role(member.Role) == role {
			return nil
		}
	}

	if Role(member.Role) == RoleOwner && slices.Contains(requiredRoles, RoleAdmin) {
		return nil
	}

	return fmt.Errorf("insufficient permissions: required %v, have %s", requiredRoles, member.Role)
}
