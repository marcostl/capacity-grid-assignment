package capacity

import (
	"context"
	"reflect"
	"testing"
	"time"
)

type repositoryStub struct {
	people   []Person
	gotRange Range
	gotLimit int
}

func (r *repositoryStub) ListPeopleCapacity(_ context.Context, dateRange Range, _ int64, limit int) ([]Person, error) {
	r.gotRange = dateRange
	r.gotLimit = limit
	return r.people, nil
}

func TestServiceNormalizesWeeksAndPaginates(t *testing.T) {
	t.Parallel()

	repository := &repositoryStub{people: []Person{{ID: 10}, {ID: 20}, {ID: 30}}}
	service := NewService(repository)
	response, err := service.Get(context.Background(), Request{
		From:  mustDate(t, "2026-01-07"),
		To:    mustDate(t, "2026-01-13"),
		Limit: 2,
	})
	if err != nil {
		t.Fatal(err)
	}

	wantRange := Range{From: mustDate(t, "2026-01-05"), To: mustDate(t, "2026-01-18")}
	if !reflect.DeepEqual(repository.gotRange, wantRange) {
		t.Fatalf("repository range = %#v, want %#v", repository.gotRange, wantRange)
	}
	if repository.gotLimit != 3 {
		t.Fatalf("repository limit = %d, want 3", repository.gotLimit)
	}
	if len(response.Weeks) != 2 || response.Weeks[0].Start != wantRange.From || response.Weeks[1].End != wantRange.To {
		t.Fatalf("weeks = %#v, want two complete weeks", response.Weeks)
	}
	if len(response.People) != 2 {
		t.Fatalf("people count = %d, want 2", len(response.People))
	}
	if response.Pagination.NextCursor == nil || *response.Pagination.NextCursor != 20 {
		t.Fatalf("next cursor = %v, want 20", response.Pagination.NextCursor)
	}
}

func TestServiceRejectsInvalidRange(t *testing.T) {
	t.Parallel()

	service := NewService(&repositoryStub{})
	_, err := service.Get(context.Background(), Request{
		From: mustDate(t, "2026-02-01"),
		To:   mustDate(t, "2026-01-01"),
	})
	if err == nil || err.Error() != "The ‘to’ date must be on or after the ‘from’ date." {
		t.Fatalf("error = %v", err)
	}
}

func mustDate(t *testing.T, value string) time.Time {
	t.Helper()
	date, err := time.Parse(time.DateOnly, value)
	if err != nil {
		t.Fatal(err)
	}
	return date
}
