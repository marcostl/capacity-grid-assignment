package people

import (
	"context"
	"fmt"
	"math"
)

const maxWeeklyHours = 168

type Repository interface {
	UpdateWeeklyHours(ctx context.Context, personID int64, weeklyHours float64) (Person, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) UpdateWeeklyHours(
	ctx context.Context,
	personID int64,
	weeklyHours float64,
) (Person, error) {
	if personID <= 0 {
		return Person{}, &ValidationError{Message: "The person ID is invalid."}
	}
	if math.IsNaN(weeklyHours) || math.IsInf(weeklyHours, 0) || weeklyHours < 0 {
		return Person{}, &ValidationError{Message: "Weekly hours must be zero or greater."}
	}
	if weeklyHours > maxWeeklyHours {
		return Person{}, &ValidationError{Message: "Weekly hours cannot exceed 168."}
	}

	person, err := s.repository.UpdateWeeklyHours(ctx, personID, weeklyHours)
	if err != nil {
		return Person{}, fmt.Errorf("update weekly hours: %w", err)
	}
	return person, nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
