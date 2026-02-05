package users

import (
	"encoding/json"
	"net/http"

	repo "github.com/envm-org/envm/internal/adapters/postgresql/sqlc"
	HTTPwriter "github.com/envm-org/envm/pkg/HTTPwriter"
	"github.com/jackc/pgx/v5/pgtype"
)

type handler struct {
	service    Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service:    service,
	}
}

func (h *handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, users)
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, user)
}

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	var userID pgtype.UUID
	if err := userID.Scan(id); err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, user)
}

func (h *handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, user)
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	var userID pgtype.UUID
	if err := userID.Scan(id); err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	var tempUser repo.UpdateUserParams
	if err := json.NewDecoder(r.Body).Decode(&tempUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tempUser.ID = userID

	user, err := h.service.UpdateUser(r.Context(), tempUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, user)
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	var userID pgtype.UUID
	if err := userID.Scan(id); err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	HTTPwriter.JSON(w, http.StatusOK, nil)
}