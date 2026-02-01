package users

import (
	"net/http"

	HTTPwriter "github.com/envm-org/envm/pkg/HTTPwriter"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	HTTPwriter.JSON(w, http.StatusOK, nil)
}

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) {
	HTTPwriter.JSON(w, http.StatusOK, nil)
}

func (h *handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	HTTPwriter.JSON(w, http.StatusOK, nil)
}

func (h *handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	HTTPwriter.JSON(w, http.StatusOK, nil)
}