package capacity

import (
	"context"
	"fmt"
	"time"
)

// Repository is the outbound port used to retrieve capacity projections.
type Repository interface {
	ListPeopleCapacity(ctx context.Context, dateRange Range, cursor int64, limit int) ([]Person, error)
}

// Service implements the capacity use case independently of HTTP and storage.
type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Get(ctx context.Context, request Request) (Response, error) {
	if request.To.Before(request.From) {
		return Response{}, validationError("The ‘to’ date must be on or after the ‘from’ date.")
	}
	if request.To.After(request.From.AddDate(2, 0, 0)) {
		return Response{}, validationError("Choose a date range of two years or less.")
	}
	if request.Limit == 0 {
		request.Limit = DefaultPageSize
	}
	if request.Limit < 0 {
		return Response{}, validationError("The page size must be a positive whole number.")
	}
	if request.Limit > MaxPageSize {
		return Response{}, validationError(fmt.Sprintf("The page size cannot exceed %d people.", MaxPageSize))
	}
	if request.Cursor < 0 {
		return Response{}, validationError("The pagination cursor is invalid.")
	}

	dateRange := Range{
		From: startOfWeek(request.From),
		To:   startOfWeek(request.To).AddDate(0, 0, 6),
	}
	people, err := s.repository.ListPeopleCapacity(ctx, dateRange, request.Cursor, request.Limit+1)
	if err != nil {
		return Response{}, fmt.Errorf("list people capacity: %w", err)
	}

	pagination := Pagination{Limit: request.Limit}
	if len(people) > request.Limit {
		nextCursor := people[request.Limit-1].ID
		pagination.NextCursor = &nextCursor
		people = people[:request.Limit]
	}

	return Response{
		Range:      dateRange,
		Weeks:      makeWeeks(dateRange),
		People:     people,
		Pagination: pagination,
	}, nil
}

func validationError(message string) error {
	return &ValidationError{Message: message}
}

func startOfWeek(date time.Time) time.Time {
	daysSinceMonday := (int(date.Weekday()) + 6) % 7
	return date.AddDate(0, 0, -daysSinceMonday)
}

func makeWeeks(dateRange Range) []Week {
	weeks := make([]Week, 0, int(dateRange.To.Sub(dateRange.From).Hours()/(24*7))+1)
	for start := dateRange.From; !start.After(dateRange.To); start = start.AddDate(0, 0, 7) {
		weeks = append(weeks, Week{Start: start, End: start.AddDate(0, 0, 6)})
	}
	return weeks
}
