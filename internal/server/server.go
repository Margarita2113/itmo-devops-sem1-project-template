package server

import (
	"fmt"
	"net/http"
	"project_sem/internal/postgres"
)

type Handler struct {
	db postgres.Postgres
}

func NewServer(db postgres.Postgres) error {
	h := Handler{db: db}

	mux := http.NewServeMux()
	mux.HandleFunc(`/api/v0/prices`, h.Handler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		return fmt.Errorf("error server %w", err)
	}
	return nil
}

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		h.POSTHandler(w, r)
	case http.MethodGet:
		h.GetHandler(w, r)
	default:
		http.Error(w, "Only POST or GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
}
