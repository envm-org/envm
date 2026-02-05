package project

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

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var userID pgtype.UUID
	userID.Scan(claims.UserID)

	if err := h.authorizer.HasRole(r.Context(), userID, organizationID, auth.RoleOwner, auth.RoleAdmin, auth.RoleMember); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var projects []repo.Project
	var err error

	if authErr := h.authorizer.HasRole(r.Context(), userID, organizationID, auth.RoleOwner, auth.RoleAdmin); authErr == nil {
		projects, err = h.service.ListProjects(r.Context(), organizationID)
	} else {
		projects, err = h.service.ListProjectsForMember(r.Context(), organizationID, userID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, projects)
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

	project, err := h.service.GetProject(r.Context(), req.ProjectID)
	if err != nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var userID pgtype.UUID
	userID.Scan(claims.UserID)

	if err := h.authorizer.HasRole(r.Context(), userID, project.OrganizationID, auth.RoleOwner, auth.RoleAdmin); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
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

	project, err := h.service.GetProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var userID pgtype.UUID
	userID.Scan(claims.UserID)

	if err := h.authorizer.HasRole(r.Context(), userID, project.OrganizationID, auth.RoleOwner, auth.RoleAdmin); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
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

	project, err := h.service.GetProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var userID pgtype.UUID
	userID.Scan(claims.UserID)

	if err := h.authorizer.HasRole(r.Context(), userID, project.OrganizationID, auth.RoleOwner, auth.RoleAdmin); err != nil {
		if _, err := h.service.GetMember(r.Context(), projectID, userID); err != nil {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
	}

	members, err := h.service.ListMembers(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, members)
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
	userID.Scan(claims.UserID)

	if err := h.authorizer.HasRole(r.Context(), userID, tempProject.OrganizationID, auth.RoleOwner, auth.RoleAdmin); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
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

	claims, ok := r.Context().Value(middleware.UserKey).(*authPkg.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var userID pgtype.UUID
	userID.Scan(claims.UserID)

	if err := h.authorizer.HasRole(r.Context(), userID, project.OrganizationID, auth.RoleOwner, auth.RoleAdmin, auth.RoleMember); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
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

	
	existingProject, err := h.service.GetProject(r.Context(), projectID)
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
	userID.Scan(claims.UserID)

	if err := h.authorizer.HasRole(r.Context(), userID, existingProject.OrganizationID, auth.RoleOwner, auth.RoleAdmin); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
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

	existingProject, err := h.service.GetProject(r.Context(), projectID)
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
	userID.Scan(claims.UserID)

	if err := h.authorizer.HasRole(r.Context(), userID, existingProject.OrganizationID, auth.RoleOwner, auth.RoleAdmin); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = h.service.DeleteProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, nil)
}
