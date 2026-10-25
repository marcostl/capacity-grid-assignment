package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"capacity/api/internal/capacity"
)

// CapacityService is the inbound port required by CapacityHandler.
type CapacityService interface {
	Get(ctx context.Context, request capacity.Request) (capacity.Response, error)
}

type CapacityHandler struct {
	service CapacityService
}

func NewCapacityHandler(service CapacityService) *CapacityHandler {
	return &CapacityHandler{service: service}
}

func (h *CapacityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request, err := parseCapacityRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.service.Get(r.Context(), request)
	if err != nil {
		var validationErr *capacity.ValidationError
		if errors.As(err, &validationErr) {
			writeError(w, http.StatusBadRequest, validationErr.Message)
			return
		}
		log.Printf("get capacity: %v", err)
		writeError(w, http.StatusInternalServerError, "Capacity could not be loaded. Please try again.")
		return
	}

	writeJSON(w, http.StatusOK, newCapacityResponse(response))
}

func parseCapacityRequest(r *http.Request) (capacity.Request, error) {
	from, err := parseRequiredDate(r, "from")
	if err != nil {
		return capacity.Request{}, err
	}
	to, err := parseRequiredDate(r, "to")
	if err != nil {
		return capacity.Request{}, err
	}
	limit, err := parseOptionalPositiveInt(r, "limit")
	if err != nil {
		return capacity.Request{}, &capacity.ValidationError{Message: "The page size must be a positive whole number."}
	}
	cursor, err := parseOptionalNonNegativeInt(r, "cursor")
	if err != nil {
		return capacity.Request{}, &capacity.ValidationError{Message: "The pagination cursor is invalid."}
	}
	return capacity.Request{From: from, To: to, Limit: limit, Cursor: cursor}, nil
}

func parseRequiredDate(r *http.Request, name string) (time.Time, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return time.Time{}, &capacity.ValidationError{Message: "Choose a ‘" + name + "’ date."}
	}
	parsed, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return time.Time{}, &capacity.ValidationError{Message: "The ‘" + name + "’ date must use the format YYYY-MM-DD."}
	}
	return parsed, nil
}

func parseOptionalPositiveInt(r *http.Request, name string) (int, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return 0, errors.New("invalid positive integer")
	}
	return parsed, nil
}

func parseOptionalNonNegativeInt(r *http.Request, name string) (int64, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed < 0 {
		return 0, errors.New("invalid non-negative integer")
	}
	return parsed, nil
}

type capacityResponse struct {
	Range      capacityRange      `json:"range"`
	Weeks      []capacityWeek     `json:"weeks"`
	People     []personCapacity   `json:"people"`
	Pagination capacityPagination `json:"pagination"`
}

type capacityRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type capacityWeek struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type personCapacity struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	WeeklyHours    float64   `json:"weeklyHours"`
	AllocatedHours []float64 `json:"allocatedHours"`
}

type capacityPagination struct {
	Limit      int    `json:"limit"`
	NextCursor string `json:"nextCursor,omitempty"`
}

func newCapacityResponse(response capacity.Response) capacityResponse {
	weeks := make([]capacityWeek, len(response.Weeks))
	for i, week := range response.Weeks {
		weeks[i] = capacityWeek{
			Start: week.Start.Format(time.DateOnly),
			End:   week.End.Format(time.DateOnly),
		}
	}
	people := make([]personCapacity, len(response.People))
	for i, person := range response.People {
		people[i] = personCapacity{
			ID:             person.ID,
			Name:           person.Name,
			WeeklyHours:    person.WeeklyHours,
			AllocatedHours: person.AllocatedHours,
		}
	}
	pagination := capacityPagination{Limit: response.Pagination.Limit}
	if response.Pagination.NextCursor != nil {
		pagination.NextCursor = strconv.FormatInt(*response.Pagination.NextCursor, 10)
	}
	return capacityResponse{
		Range: capacityRange{
			From: response.Range.From.Format(time.DateOnly),
			To:   response.Range.To.Format(time.DateOnly),
		},
		Weeks:      weeks,
		People:     people,
		Pagination: pagination,
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
