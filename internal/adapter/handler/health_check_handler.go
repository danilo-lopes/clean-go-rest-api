package handler

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

type HealthCheckHandler struct {
	db *sql.DB
}

func NewHealthCheckHandler(db *sql.DB) *HealthCheckHandler {
	return &HealthCheckHandler{db: db}
}

func (h *HealthCheckHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/healthz", h.HealthCheck).Methods("GET")
	r.HandleFunc("/readyz", h.ReadyCheck).Methods("GET")
}

func (h *HealthCheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (h *HealthCheckHandler) ReadyCheck(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("db unavailable"))
		return
	}
	if err := h.db.Ping(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("db unavailable"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}
