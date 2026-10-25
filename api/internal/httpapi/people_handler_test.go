package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"capacity/api/internal/people"
)

type peopleServiceStub struct {
	person people.Person
	err    error
	gotID  int64
	got    float64
}

func (s *peopleServiceStub) UpdateWeeklyHours(_ context.Context, id int64, hours float64) (people.Person, error) {
	s.gotID = id
	s.got = hours
	return s.person, s.err
}

func TestPeopleHandlerUpdatesWeeklyHours(t *testing.T) {
	t.Parallel()

	service := &peopleServiceStub{person: people.Person{ID: 42, Name: "Vigdis", WeeklyHours: 32}}
	recorder := httptest.NewRecorder()
	request := newPatchRequest("42", `{"weeklyHours":32}`)

	NewPeopleHandler(service).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if service.gotID != 42 || service.got != 32 {
		t.Fatalf("service called with %d, %v", service.gotID, service.got)
	}
	var response personResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if response.ETag == "" || recorder.Header().Get("ETag") != response.ETag {
		t.Fatalf("response ETag = %q, header = %q", response.ETag, recorder.Header().Get("ETag"))
	}
}

func TestPeopleHandlerRejectsInvalidBody(t *testing.T) {
	t.Parallel()

	for _, body := range []string{`{}`, `{"weeklyHours":"many"}`, `{"weeklyHours":40,"other":1}`, `{"weeklyHours":40}{}`} {
		recorder := httptest.NewRecorder()
		NewPeopleHandler(&peopleServiceStub{}).ServeHTTP(recorder, newPatchRequest("1", body))
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("body %q: status = %d, want %d", body, recorder.Code, http.StatusBadRequest)
		}
	}
}

func TestPeopleHandlerMapsServiceErrors(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name   string
		err    error
		status int
		want   string
	}{
		{name: "not found", err: people.ErrNotFound, status: http.StatusNotFound, want: "That person could not be found."},
		{name: "internal", err: errors.New("database unavailable"), status: http.StatusInternalServerError, want: "Weekly hours could not be updated. Please try again."},
	} {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			service := &peopleServiceStub{err: test.err}
			NewPeopleHandler(service).ServeHTTP(recorder, newPatchRequest("1", `{"weeklyHours":40}`))
			if recorder.Code != test.status {
				t.Fatalf("status = %d, want %d", recorder.Code, test.status)
			}
			var response map[string]string
			if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
				t.Fatal(err)
			}
			if response["error"] != test.want {
				t.Fatalf("error = %q, want %q", response["error"], test.want)
			}
		})
	}
}

func newPatchRequest(id, body string) *http.Request {
	request := httptest.NewRequest(http.MethodPatch, "/api/people/"+id, strings.NewReader(body))
	request.SetPathValue("id", id)
	return request
}
