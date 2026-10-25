package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"capacity/api/internal/capacity"
)

type capacityServiceStub struct {
	err error
}

func (s capacityServiceStub) Get(_ context.Context, _ capacity.Request) (capacity.Response, error) {
	return capacity.Response{}, s.err
}

func TestCapacityHandlerHidesInternalErrors(t *testing.T) {
	t.Parallel()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/capacity?from=2026-01-01&to=2026-01-02", nil)
	handler := NewCapacityHandler(capacityServiceStub{err: errors.New("database connection failed")})

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	var response map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	const want = "Capacity could not be loaded. Please try again."
	if response["error"] != want {
		t.Fatalf("error = %q, want %q", response["error"], want)
	}
}

func TestCapacityHandlerRejectsInvalidHTTPInput(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name    string
		query   string
		message string
	}{
		{name: "missing from", query: "to=2026-01-01", message: "Choose a ‘from’ date."},
		{name: "malformed date", query: "from=not-a-date&to=2026-01-01", message: "The ‘from’ date must use the format YYYY-MM-DD."},
		{name: "invalid limit", query: "from=2026-01-01&to=2026-01-02&limit=0", message: "The page size must be a positive whole number."},
		{name: "invalid cursor", query: "from=2026-01-01&to=2026-01-02&cursor=next", message: "The pagination cursor is invalid."},
	} {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/api/capacity?"+test.query, nil)

			NewCapacityHandler(capacityServiceStub{}).ServeHTTP(recorder, request)

			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
			}
			var response map[string]string
			if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if response["error"] != test.message {
				t.Fatalf("error = %q, want %q", response["error"], test.message)
			}
		})
	}
}
