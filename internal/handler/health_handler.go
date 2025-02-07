package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type HealthHandler struct {
	isReady *atomic.Bool
	version string
	db      *mongo.Database
}

type healthResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version,omitempty"`
	CheckTime string            `json:"check_time"`
	Checks    map[string]string `json:"checks,omitempty"`
}

type readinessResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

func NewHealthHandler(version string, db *mongo.Database) *HealthHandler {
	ready := &atomic.Bool{}
	ready.Store(false)

	return &HealthHandler{
		isReady: ready,
		version: version,
		db:      db,
	}
}

func (h *HealthHandler) MarkAsReady() {
	h.isReady.Store(true)
}

func (h *HealthHandler) MarkAsNotReady() {
	h.isReady.Store(false)
}

func (h *HealthHandler) checkMongoDB() string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := h.db.Client().Ping(ctx, readpref.Primary()); err != nil {
		return "DOWN"
	}
	return "UP"
}

func (h *HealthHandler) LivenessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	checks := map[string]string{
		"mongodb": h.checkMongoDB(),
	}

	status := "UP"
	for _, check := range checks {
		if check != "UP" {
			status = "DOWN"
			break
		}
	}

	response := healthResponse{
		Status:    status,
		Version:   h.version,
		CheckTime: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	}

	if status != "UP" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(response)
}

func (h *HealthHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	checks := map[string]string{
		"mongodb": h.checkMongoDB(),
	}

	if !h.isReady.Load() {
		checks["ready"] = "NOT_READY"
	} else {
		checks["ready"] = "READY"
	}

	status := "READY"
	for _, check := range checks {
		if check != "UP" && check != "READY" {
			status = "NOT_READY"
			break
		}
	}

	if status != "READY" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	response := readinessResponse{
		Status:    status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	}
	json.NewEncoder(w).Encode(response)
}
