package postgres

import (
	"context"
	"errors"
	"fmt"

	"capacity/api/internal/people"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PeopleRepository struct {
	db *pgxpool.Pool
}

func NewPeopleRepository(db *pgxpool.Pool) *PeopleRepository {
	return &PeopleRepository{db: db}
}

func (r *PeopleRepository) UpdateWeeklyHours(
	ctx context.Context,
	personID int64,
	weeklyHours float64,
) (people.Person, error) {
	var person people.Person
	err := r.db.QueryRow(ctx, `
		UPDATE people
		SET weekly_hours = $2
		WHERE id = $1
		RETURNING id, name, weekly_hours::double precision
	`, personID, weeklyHours).Scan(&person.ID, &person.Name, &person.WeeklyHours)
	if errors.Is(err, pgx.ErrNoRows) {
		return people.Person{}, people.ErrNotFound
	}
	if err != nil {
		return people.Person{}, fmt.Errorf("update person: %w", err)
	}
	return person, nil
}
