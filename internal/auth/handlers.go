package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	HTTPwriter "github.com/envm-org/envm/pkg/HTTPwriter"
	"github.com/envm-org/envm/pkg/auth"
	"github.com/envm-org/envm/pkg/env"
	"github.com/envm-org/envm/pkg/validator"
	goValidator "github.com/go-playground/validator/v10"
)

type handler struct {
	service    Service
	tokenMaker auth.TokenMaker
	validate   *goValidator.Validate
}

func NewHandler(service Service, tokenMaker auth.TokenMaker) *handler {
	return &handler{
		service:    service,
		tokenMaker: tokenMaker,
		validate:   validator.New(),
	}
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	accessToken, err := h.tokenMaker.CreateToken(user.ID, user.Email, 15*time.Minute)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.service.CreateSession(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   env.GetString("ENV", "development") == "production",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   15 * 60, // 15 minutes
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   env.GetString("ENV", "development") == "production",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60, // 7 days
	})

	HTTPwriter.JSON(w, http.StatusOK, map[string]interface{}{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "missing refresh token", http.StatusUnauthorized)
		return
	}

	user, err := h.service.ValidateRefreshToken(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	accessToken, err := h.tokenMaker.CreateToken(user.ID, user.Email, 15*time.Minute)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   env.GetString("ENV", "development") == "production",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   15 * 60,
	})

	HTTPwriter.JSON(w, http.StatusOK, map[string]interface{}{
		"access_token": accessToken,
	})
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "missing refresh token", http.StatusBadRequest)
		return
	}

	if err := h.service.Logout(r.Context(), cookie.Value); err != nil {
		// Log error but don't block logout
		fmt.Printf("Logout Error: %v\n", err)
	}

	// Clear cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
	})

	HTTPwriter.JSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	HTTPwriter.JSON(w, http.StatusCreated, user)
}

func (h *handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.ForgotPassword(r.Context(), req.Email); err != nil {
		fmt.Printf("Forgot Password Error: %v\n", err)
	}

	HTTPwriter.JSON(w, http.StatusOK, map[string]string{"message": "If email exists, instructions have been sent"})
}

func (h *handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token       string `json:"token" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=8"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.ResetPassword(r.Context(), req.Token, req.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	HTTPwriter.JSON(w, http.StatusOK, map[string]string{"message": "Password updated successfully"})
}
