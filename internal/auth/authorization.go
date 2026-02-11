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
	HasProjectRole(ctx context.Context, userID, projectID pgtype.UUID, requiredRoles ...Role) error
}

type authorizer struct {
	repo *repo.Queries
}

func NewAuthorizer(repo *repo.Queries) Authorizer {
	return &authorizer{repo: repo}
}

func (a *authorizer) HasProjectRole(ctx context.Context, userID, projectID pgtype.UUID, requiredRoles ...Role) error {
	member, err := a.repo.GetProjectMember(ctx, repo.GetProjectMemberParams{
		ProjectID: projectID,
		UserID:    userID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user is not a member of this project")
		}
		return fmt.Errorf("failed to check membership: %w", err)
	}

	for _, role := range requiredRoles {
		if Role(member.Role) == role {
			return nil
		}
	}

	userRole := Role(member.Role)

	// fast path exact match
	if slices.Contains(requiredRoles, userRole) {
		return nil
	}

	// Hierarchy: Owner > Admin > Member
	if userRole == RoleOwner {
		return nil // Owner has all permissions
	}
	if userRole == RoleAdmin && slices.Contains(requiredRoles, RoleMember) {
		return nil // Admin has Member permissions
	}

	return fmt.Errorf("insufficient permissions: required %v, have %s", requiredRoles, member.Role)
}
