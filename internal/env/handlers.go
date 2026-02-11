package env

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

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var userID pgtype.UUID
	if err := userID.Scan(claims.UserID); err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	// Check if member
	if err := h.authorizer.HasProjectRole(r.Context(), userID, projectID, auth.RoleMember, auth.RoleAdmin, auth.RoleOwner); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
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

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var userID pgtype.UUID
	if err := userID.Scan(claims.UserID); err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	// Check if admin/owner
	if err := h.authorizer.HasProjectRole(r.Context(), userID, tempEnv.ProjectID, auth.RoleAdmin, auth.RoleOwner); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
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

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var userID pgtype.UUID
	if err := userID.Scan(claims.UserID); err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	if err := h.authorizer.HasProjectRole(r.Context(), userID, env.ProjectID, auth.RoleMember, auth.RoleAdmin, auth.RoleOwner); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
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

	// Check existing env to get ProjectID
	existingEnv, err := h.service.GetEnv(r.Context(), envID)
	if err != nil {
		http.Error(w, "environment not found", http.StatusNotFound)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var userID pgtype.UUID
	if err := userID.Scan(claims.UserID); err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	// Check if admin/owner
	if err := h.authorizer.HasProjectRole(r.Context(), userID, existingEnv.ProjectID, auth.RoleAdmin, auth.RoleOwner); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

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

	// Check existing env to get ProjectID
	existingEnv, err := h.service.GetEnv(r.Context(), envID)
	if err != nil {
		http.Error(w, "environment not found", http.StatusNotFound)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var userID pgtype.UUID
	if err := userID.Scan(claims.UserID); err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	// Check if admin/owner
	if err := h.authorizer.HasProjectRole(r.Context(), userID, existingEnv.ProjectID, auth.RoleAdmin, auth.RoleOwner); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = h.service.DeleteEnv(r.Context(), envID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, nil)
}

func (h *handler) CreateVariable(w http.ResponseWriter, r *http.Request) {
	var params repo.CreateVariableParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Need to get ProjectID from EnvironmentID
	env, err := h.service.GetEnv(r.Context(), params.EnvironmentID)
	if err != nil {
		http.Error(w, "environment not found", http.StatusNotFound)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var userID pgtype.UUID
	if err := userID.Scan(claims.UserID); err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	// Check if admin/owner
	if err := h.authorizer.HasProjectRole(r.Context(), userID, env.ProjectID, auth.RoleAdmin, auth.RoleOwner); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	variable, err := h.service.CreateVariable(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, variable)
}

func (h *handler) UpdateVariable(w http.ResponseWriter, r *http.Request) {
	var params repo.UpdateVariableParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Need to get ProjectID from EnvironmentID
	env, err := h.service.GetEnv(r.Context(), params.EnvironmentID)
	if err != nil {
		http.Error(w, "environment not found", http.StatusNotFound)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var userID pgtype.UUID
	if err := userID.Scan(claims.UserID); err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	// Check if admin/owner
	if err := h.authorizer.HasProjectRole(r.Context(), userID, env.ProjectID, auth.RoleAdmin, auth.RoleOwner); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	variable, err := h.service.UpdateVariable(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, variable)
}

func (h *handler) DeleteVariable(w http.ResponseWriter, r *http.Request) {
	environmentIDStr := r.URL.Query().Get("environment_id")
	if environmentIDStr == "" {
		http.Error(w, "environment_id is required", http.StatusBadRequest)
		return
	}
	var environmentID pgtype.UUID
	if err := environmentID.Scan(environmentIDStr); err != nil {
		http.Error(w, "invalid environment_id format", http.StatusBadRequest)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}

	// Need to get ProjectID from EnvironmentID
	env, err := h.service.GetEnv(r.Context(), environmentID)
	if err != nil {
		http.Error(w, "environment not found", http.StatusNotFound)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var userID pgtype.UUID
	if err := userID.Scan(claims.UserID); err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	// Check if admin/owner
	if err := h.authorizer.HasProjectRole(r.Context(), userID, env.ProjectID, auth.RoleAdmin, auth.RoleOwner); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = h.service.DeleteVariable(r.Context(), environmentID, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, nil)
}

func (h *handler) ListVariables(w http.ResponseWriter, r *http.Request) {
	environmentIDStr := r.URL.Query().Get("environment_id")
	if environmentIDStr == "" {
		http.Error(w, "environment_id is required", http.StatusBadRequest)
		return
	}
	var environmentID pgtype.UUID
	if err := environmentID.Scan(environmentIDStr); err != nil {
		http.Error(w, "invalid environment_id format", http.StatusBadRequest)
		return
	}

	// Need to get ProjectID from EnvironmentID
	env, err := h.service.GetEnv(r.Context(), environmentID)
	if err != nil {
		http.Error(w, "environment not found", http.StatusNotFound)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var userID pgtype.UUID
	if err := userID.Scan(claims.UserID); err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	// Check if member
	if err := h.authorizer.HasProjectRole(r.Context(), userID, env.ProjectID, auth.RoleMember, auth.RoleAdmin, auth.RoleOwner); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	vars, err := h.service.ListVariables(r.Context(), environmentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, vars)
}
