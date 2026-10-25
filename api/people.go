package main

import (
	"net/http"

	"capacity/api/internal/httpapi"
	"capacity/api/internal/people"
	postgresadapter "capacity/api/internal/postgres"
)

// handleUpdatePerson is the composition root for updating a person's capacity.
func (s *server) handleUpdatePerson(w http.ResponseWriter, r *http.Request) {
	repository := postgresadapter.NewPeopleRepository(s.db)
	service := people.NewService(repository)
	handler := httpapi.NewPeopleHandler(service)
	handler.ServeHTTP(w, r)
}
