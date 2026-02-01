package org

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

func (h *handler) CreateOrg(w http.ResponseWriter, r *http.Request) {
	HTTPwriter.JSON(w, http.StatusOK, nil)
}

func (h *handler) GetOrg(w http.ResponseWriter, r *http.Request) {
	HTTPwriter.JSON(w, http.StatusOK, nil)
}

func (h *handler) ListOrgs(w http.ResponseWriter, r *http.Request) {
	HTTPwriter.JSON(w, http.StatusOK, nil)
}