package main

import (
	"net/http"

	"capacity/api/internal/capacity"
	"capacity/api/internal/httpapi"
	postgresadapter "capacity/api/internal/postgres"
)

// handleCapacity is the composition root for the capacity feature. Keeping the
// feature wiring here lets main.go remain concerned only with server startup and
// route registration.
func (s *server) handleCapacity(w http.ResponseWriter, r *http.Request) {
	repository := postgresadapter.NewCapacityRepository(s.db)
	service := capacity.NewService(repository)
	handler := httpapi.NewCapacityHandler(service)
	handler.ServeHTTP(w, r)
}
