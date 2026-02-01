package env

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

func (h *handler) ListEnvs(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	if projectIDStr == "" {
		http.Error(w, "project_id is required", http.StatusBadRequest)
		return
	}
	var projectID pgtype.UUID
	if err := projectID.Scan(projectIDStr); err != nil {
		http.Error(w, "invalid project_id format", http.StatusBadRequest)
		return
	}

	envs, err := h.service.ListEnvs(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, envs)
}

func (h *handler) CreateEnv(w http.ResponseWriter, r *http.Request) {
	var tempEnv repo.CreateEnvironmentParams
	if err := json.NewDecoder(r.Body).Decode(&tempEnv); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	env, err := h.service.CreateEnv(r.Context(), tempEnv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, env)
}

func (h *handler) GetEnv(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	var envID pgtype.UUID
	if err := envID.Scan(id); err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	env, err := h.service.GetEnv(r.Context(), envID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, env)
}

func (h *handler) UpdateEnv(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	var envID pgtype.UUID
	if err := envID.Scan(id); err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	var tempEnv repo.UpdateEnvironmentParams
	if err := json.NewDecoder(r.Body).Decode(&tempEnv); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tempEnv.ID = envID

	env, err := h.service.UpdateEnv(r.Context(), tempEnv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, env)
}

func (h *handler) DeleteEnv(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	var envID pgtype.UUID
	if err := envID.Scan(id); err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteEnv(r.Context(), envID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, nil)
}
