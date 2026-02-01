package org

import (
	"encoding/json"
	"net/http"

	repo "github.com/envm-org/envm/internal/adapters/postgresql/sqlc"
	HTTPwriter "github.com/envm-org/envm/pkg/HTTPwriter"
	"github.com/jackc/pgx/v5/pgtype"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
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

	err := h.service.DeleteOrg(r.Context(), orgID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, nil)
}