package people

import (
	"context"
	"errors"
	"testing"
)

type repositoryStub struct {
	person         Person
	err            error
	gotID          int64
	gotWeeklyHours float64
}

func (r *repositoryStub) UpdateWeeklyHours(_ context.Context, id int64, hours float64) (Person, error) {
	r.gotID = id
	r.gotWeeklyHours = hours
	return r.person, r.err
}

func TestUpdateWeeklyHours(t *testing.T) {
	t.Parallel()

	repository := &repositoryStub{person: Person{ID: 7, Name: "Sanne", WeeklyHours: 32.5}}
	person, err := NewService(repository).UpdateWeeklyHours(context.Background(), 7, 32.5)
	if err != nil {
		t.Fatal(err)
	}
	if repository.gotID != 7 || repository.gotWeeklyHours != 32.5 {
		t.Fatalf("repository called with %d, %v", repository.gotID, repository.gotWeeklyHours)
	}
	if person != repository.person {
		t.Fatalf("person = %#v, want %#v", person, repository.person)
	}
}

func TestUpdateWeeklyHoursValidation(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name  string
		id    int64
		hours float64
	}{
		{name: "invalid ID", id: 0, hours: 40},
		{name: "negative hours", id: 1, hours: -1},
		{name: "too many hours", id: 1, hours: 169},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewService(&repositoryStub{}).UpdateWeeklyHours(context.Background(), test.id, test.hours)
			var validationErr *ValidationError
			if !errors.As(err, &validationErr) {
				t.Fatalf("error = %v, want ValidationError", err)
			}
		})
	}
}
