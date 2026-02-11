package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	repo "github.com/envm-org/envm/internal/adapters/postgresql/sqlc"
	"github.com/envm-org/envm/pkg/auth"
	"github.com/envm-org/envm/pkg/email"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type RegisterParams struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
}

type Service interface {
	Login(ctx context.Context, email, password string) (repo.User, error)
	Register(ctx context.Context, params RegisterParams) (repo.User, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	CreateSession(ctx context.Context, userID pgtype.UUID) (string, error)
	ValidateRefreshToken(ctx context.Context, token string) (repo.User, error)
	Logout(ctx context.Context, token string) error
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

func (s *svc) Login(ctx context.Context, email, password string) (repo.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return repo.User{}, fmt.Errorf("invalid credentials")
		}
		return repo.User{}, fmt.Errorf("failed to fetch user: %w", err)
	}

	err = auth.CheckPassword(password, user.PasswordHash)
	if err != nil {
		return repo.User{}, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

func (s *svc) Register(ctx context.Context, params RegisterParams) (repo.User, error) {
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		return repo.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	return s.repo.CreateUser(ctx, repo.CreateUserParams{
		Email:        params.Email,
		PasswordHash: hashedPassword,
		FullName:     params.FullName,
	})
}

func (s *svc) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return err
	}
	token := hex.EncodeToString(tokenBytes)

	expiresAt := time.Now().Add(1 * time.Hour) // 1 hour expiry

	err = s.repo.SetPasswordResetToken(ctx, repo.SetPasswordResetTokenParams{
		Email:                  user.Email,
		PasswordResetToken:     pgtype.Text{String: token, Valid: true},
		PasswordResetExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return err
	}

	// Send Email
	subject := "Reset your password"
	body := fmt.Sprintf("Use this token to reset your password: %s", token)
	return s.mailer.SendEmail(user.Email, subject, body)
}

func (s *svc) ResetPassword(ctx context.Context, token, newPassword string) error {
	user, err := s.repo.GetUserByResetToken(ctx, pgtype.Text{String: token, Valid: true})
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("invalid or expired token")
		}
		return err
	}

	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, repo.UpdatePasswordParams{
		ID:           user.ID,
		PasswordHash: hashedPassword,
	})
}

func (s *svc) CreateSession(ctx context.Context, userID pgtype.UUID) (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days

	_, err := s.repo.CreateRefreshToken(ctx, repo.CreateRefreshTokenParams{
		UserID:    userID,
		Token:     token,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	return token, nil
}

func (s *svc) ValidateRefreshToken(ctx context.Context, token string) (repo.User, error) {
	refreshToken, err := s.repo.GetRefreshToken(ctx, token)
	if err != nil {
		if err == pgx.ErrNoRows {
			return repo.User{}, fmt.Errorf("invalid refresh token")
		}
		return repo.User{}, err
	}

	if refreshToken.Revoked.Bool {
		return repo.User{}, fmt.Errorf("refresh token revoked")
	}

	if refreshToken.ExpiresAt.Time.Before(time.Now()) {
		return repo.User{}, fmt.Errorf("refresh token expired")
	}

	user, err := s.repo.GetUser(ctx, refreshToken.UserID)
	if err != nil {
		return repo.User{}, err
	}

	return user, nil
}

func (s *svc) Logout(ctx context.Context, token string) error {
	return s.repo.RevokeRefreshToken(ctx, token)
}
