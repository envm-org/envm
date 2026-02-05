package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	HTTPwriter "github.com/envm-org/envm/pkg/HTTPwriter"
	"github.com/envm-org/envm/pkg/auth"
	"github.com/envm-org/envm/pkg/env"
)

type handler struct {
	service    Service
	tokenMaker auth.TokenMaker
}

func NewHandler(service Service, tokenMaker auth.TokenMaker) *handler {
	return &handler{
		service:    service,
		tokenMaker: tokenMaker,
	}
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := h.tokenMaker.CreateToken(user.ID, user.Email, 24*time.Hour) 
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError) 
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   env.GetString("ENV", "development") == "production", 
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	HTTPwriter.JSON(w, http.StatusOK, map[string]interface{}{ 
		"user":  user,
		"token": token,
	})
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest) 
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
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.ForgotPassword(r.Context(), req.Email); err != nil {
		fmt.Printf("Forgot Password Error: %v\n", err)
	}

	HTTPwriter.JSON(w, http.StatusOK, map[string]string{"message": "If email exists, instructions have been sent"})
}

func (h *handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.ResetPassword(r.Context(), req.Token, req.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	HTTPwriter.JSON(w, http.StatusOK, map[string]string{"message": "Password updated successfully"})
}
