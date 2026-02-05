package org

import (
	"encoding/json"
	"net/http"

	repo "github.com/envm-org/envm/internal/adapters/postgresql/sqlc"
	"github.com/envm-org/envm/internal/auth"
	"github.com/envm-org/envm/internal/middleware"
	HTTPwriter "github.com/envm-org/envm/pkg/HTTPwriter"
	authPkg "github.com/envm-org/envm/pkg/auth"
	"github.com/jackc/pgx/v5/pgtype"
)

type handler struct {
	service    Service
	authorizer auth.Authorizer
}

func NewHandler(service Service, authorizer auth.Authorizer) *handler {
	return &handler{
		service:    service,
		authorizer: authorizer,
	}
}

func (h *handler) ListOrgs(w http.ResponseWriter, r *http.Request) {
	orgs, err := h.service.ListOrgs(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, orgs)
}

func (h *handler) CreateOrg(w http.ResponseWriter, r *http.Request) {
	var tempOrg repo.CreateOrganizationParams
	if err := json.NewDecoder(r.Body).Decode(&tempOrg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	org, err := h.service.CreateOrg(r.Context(), tempOrg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// TODO: Add the creator as Owner member
	HTTPwriter.JSON(w, http.StatusOK, org)
}

func (h *handler) GetOrg(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	var orgID pgtype.UUID
	if err := orgID.Scan(id); err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	
	var userID pgtype.UUID
	if err := userID.Scan(claims.UserID); err != nil {
		http.Error(w, "invalid user id in token", http.StatusUnauthorized)
		return
	}

	if err := h.authorizer.HasRole(r.Context(), userID, orgID, auth.RoleOwner, auth.RoleAdmin, auth.RoleMember); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	org, err := h.service.GetOrg(r.Context(), orgID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, org)
}

func (h *handler) UpdateOrg(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	var orgID pgtype.UUID
	if err := orgID.Scan(id); err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	
	var userID pgtype.UUID
	userID.Scan(claims.UserID)

	if err := h.authorizer.HasRole(r.Context(), userID, orgID, auth.RoleOwner, auth.RoleAdmin); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var tempOrg repo.UpdateOrganizationParams
	if err := json.NewDecoder(r.Body).Decode(&tempOrg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tempOrg.ID = orgID

	org, err := h.service.UpdateOrg(r.Context(), tempOrg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, org)
}

func (h *handler) DeleteOrg(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	var orgID pgtype.UUID
	if err := orgID.Scan(id); err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	
	var userID pgtype.UUID
	userID.Scan(claims.UserID)

	if err := h.authorizer.HasRole(r.Context(), userID, orgID, auth.RoleOwner); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err := h.service.DeleteOrg(r.Context(), orgID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, nil)
}

func (h *handler) InviteMember(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrganizationID pgtype.UUID `json:"organization_id"`
		Email          string      `json:"email"`
		Role           string      `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var userID pgtype.UUID
	userID.Scan(claims.UserID)

	if err := h.authorizer.HasRole(r.Context(), userID, req.OrganizationID, auth.RoleOwner, auth.RoleAdmin); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if err := h.service.InviteMember(r.Context(), req.OrganizationID, req.Email, req.Role, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusCreated, map[string]string{"message": "invitation sent"})
}

func (h *handler) AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var userID pgtype.UUID
	userID.Scan(claims.UserID)

	if err := h.service.AcceptInvitation(r.Context(), req.Token, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, map[string]string{"message": "invitation accepted"})
}