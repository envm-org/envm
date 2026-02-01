package env

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

func (h *handler) CreateEnv(w http.ResponseWriter, r *http.Request) {
	HTTPwriter.JSON(w, http.StatusOK, nil)
}

func (h *handler) GetEnv(w http.ResponseWriter, r *http.Request) {
	HTTPwriter.JSON(w, http.StatusOK, nil)
}

func (h *handler) UpdateEnv(w http.ResponseWriter, r *http.Request) {
	HTTPwriter.JSON(w, http.StatusOK, nil)
}

func (h *handler) DeleteEnv(w http.ResponseWriter, r *http.Request) {
	HTTPwriter.JSON(w, http.StatusOK, nil)
}
