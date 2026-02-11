package project

import (
	"encoding/json"
	"net/http"

	repo "github.com/envm-org/envm/internal/adapters/postgresql/sqlc"
	"github.com/envm-org/envm/internal/middleware"
	HTTPwriter "github.com/envm-org/envm/pkg/HTTPwriter"
	authPkg "github.com/envm-org/envm/pkg/auth"
	"github.com/jackc/pgx/v5/pgtype"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ListProjects(w http.ResponseWriter, r *http.Request) {
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

	projects, err := h.service.ListProjects(r.Context(), userID)
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

	project, err := h.service.CreateProject(r.Context(), tempProject, userID)
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

	// Permission check: Must be a member
	if _, err := h.service.GetMember(r.Context(), projectID, userID); err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
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

	// Permission check: Must be owner or admin
	member, err := h.service.GetMember(r.Context(), projectID, userID)
	if err != nil || (member.Role != "owner" && member.Role != "admin") {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

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

	// Permission check: Must be owner
	member, err := h.service.GetMember(r.Context(), projectID, userID)
	if err != nil || member.Role != "owner" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := h.service.DeleteProject(r.Context(), projectID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, nil)
}

func (h *handler) AddMember(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProjectID pgtype.UUID `json:"project_id"`
		UserID    pgtype.UUID `json:"user_id"`
		Role      string      `json:"role"`
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
	if err := userID.Scan(claims.UserID); err != nil {
		http.Error(w, "invalid user id in token", http.StatusUnauthorized)
		return
	}

	// Permission check: Must be owner or admin
	member, err := h.service.GetMember(r.Context(), req.ProjectID, userID)
	if err != nil || (member.Role != "owner" && member.Role != "admin") {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := h.service.AddMember(r.Context(), req.ProjectID, req.UserID, req.Role); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusCreated, map[string]string{"message": "member added"})
}

func (h *handler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	userIDStr := r.URL.Query().Get("user_id")

	var projectID, targetUserID pgtype.UUID
	if err := projectID.Scan(projectIDStr); err != nil {
		http.Error(w, "invalid project_id", http.StatusBadRequest)
		return
	}
	if err := targetUserID.Scan(userIDStr); err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
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

	// Permission check: Must be owner or admin
	member, err := h.service.GetMember(r.Context(), projectID, userID)
	if err != nil || (member.Role != "owner" && member.Role != "admin") {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := h.service.RemoveMember(r.Context(), projectID, targetUserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, map[string]string{"message": "member removed"})
}

func (h *handler) ListMembers(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	var projectID pgtype.UUID
	if err := projectID.Scan(projectIDStr); err != nil {
		http.Error(w, "invalid project_id", http.StatusBadRequest)
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

	// Permission check: Must be a member
	if _, err := h.service.GetMember(r.Context(), projectID, userID); err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	members, err := h.service.ListMembers(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, members)
}
