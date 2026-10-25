package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"capacity/api/internal/people"
)

const maxUpdatePersonBodyBytes = 4 << 10

type PeopleService interface {
	UpdateWeeklyHours(ctx context.Context, personID int64, weeklyHours float64) (people.Person, error)
}

type PeopleHandler struct {
	service PeopleService
}

func NewPeopleHandler(service PeopleService) *PeopleHandler {
	return &PeopleHandler{service: service}
}

func (h *PeopleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	personID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || personID <= 0 {
		writeError(w, http.StatusBadRequest, "The person ID is invalid.")
		return
	}

	request, err := decodeUpdatePersonRequest(w, r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	person, err := h.service.UpdateWeeklyHours(r.Context(), personID, *request.WeeklyHours)
	if err != nil {
		var validationErr *people.ValidationError
		switch {
		case errors.As(err, &validationErr):
			writeError(w, http.StatusBadRequest, validationErr.Message)
		case errors.Is(err, people.ErrNotFound):
			writeError(w, http.StatusNotFound, "That person could not be found.")
		default:
			log.Printf("update person: %v", err)
			writeError(w, http.StatusInternalServerError, "Weekly hours could not be updated. Please try again.")
		}
		return
	}

	response := newPersonResponse(person)
	w.Header().Set("ETag", response.ETag)
	writeJSON(w, http.StatusOK, response)
}

type updatePersonRequest struct {
	WeeklyHours *float64 `json:"weeklyHours"`
}

func decodeUpdatePersonRequest(w http.ResponseWriter, r *http.Request) (updatePersonRequest, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUpdatePersonBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var request updatePersonRequest
	if err := decoder.Decode(&request); err != nil {
		return updatePersonRequest{}, errors.New("Enter weekly hours as a JSON number.")
	}
	if request.WeeklyHours == nil {
		return updatePersonRequest{}, errors.New("Weekly hours are required.")
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return updatePersonRequest{}, errors.New("The request body must contain one JSON object.")
	}
	return request, nil
}

type personResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	WeeklyHours float64 `json:"weeklyHours"`
	ETag        string  `json:"etag"`
}

func newPersonResponse(person people.Person) personResponse {
	return personResponse{
		ID:          person.ID,
		Name:        person.Name,
		WeeklyHours: person.WeeklyHours,
		ETag:        personETag(person.ID, person.WeeklyHours),
	}
}
