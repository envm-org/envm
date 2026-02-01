package project

import (
	"encoding/json"
	"net/http"

	repo "github.com/envm-org/envm/internal/adapters/postgresql/sqlc"
	"github.com/envm-org/envm/pkg/HTTPwriter"
	"github.com/jackc/pgx/v5/pgtype"
)

type handler struct {
	service Service
}


func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	organizationIDStr := r.URL.Query().Get("organization_id")
	if organizationIDStr == "" {
		http.Error(w, "organization_id is required", http.StatusBadRequest)
		return
	}
	var organizationID pgtype.UUID
	if err := organizationID.Scan(organizationIDStr); err != nil {
		http.Error(w, "invalid organization_id format", http.StatusBadRequest)
		return
	}

	projects, err := h.service.ListProjects(r.Context(), organizationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, projects)
}

func (h *handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var tempProject repo.CreateProjectParams
	if err := json.NewDecoder(r.Body).Decode(&tempProject); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	project, err := h.service.CreateProject(r.Context(), tempProject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, project)
}


func (h *handler) GetProject(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	var projectID pgtype.UUID
	if err := projectID.Scan(id); err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}
	project, err := h.service.GetProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, project)
}



func (h *handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	var projectID pgtype.UUID
	if err := projectID.Scan(id); err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}
	var tempProject repo.UpdateProjectParams
	if err := json.NewDecoder(r.Body).Decode(&tempProject); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tempProject.ID = projectID
	project, err := h.service.UpdateProject(r.Context(), tempProject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, project)
}


func (h *handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	var projectID pgtype.UUID
	if err := projectID.Scan(id); err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}
	err := h.service.DeleteProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, nil)
}
