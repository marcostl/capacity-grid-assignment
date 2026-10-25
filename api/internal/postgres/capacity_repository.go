package postgres

import (
	"context"
	"fmt"
	"time"

	"capacity/api/internal/capacity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CapacityRepository is the PostgreSQL adapter for the capacity outbound port.
type CapacityRepository struct {
	db *pgxpool.Pool
}

func NewCapacityRepository(db *pgxpool.Pool) *CapacityRepository {
	return &CapacityRepository{db: db}
}

func (r *CapacityRepository) ListPeopleCapacity(
	ctx context.Context,
	dateRange capacity.Range,
	cursor int64,
	limit int,
) ([]capacity.Person, error) {
	// Page people first so assignment work stays bounded. The CTEs generate
	// Monday-starting weeks, sum each assignment's inclusive overlap with each
	// week, then zero-fill and aggregate those totals into one ordered array per
	// person. Calendar days, including weekends, all contribute to allocation.
	rows, err := r.db.Query(ctx, `
		WITH selected_people AS (
			SELECT id, name, weekly_hours
			FROM people
			WHERE id > $3
			ORDER BY id
			LIMIT $4
		),
		weeks AS (
			SELECT week_start::date
			FROM generate_series($1::date, $2::date, interval '7 days') AS week_start
		),
		allocations AS (
			SELECT
				p.id AS person_id,
				w.week_start,
				SUM(
					a.hours_per_day *
					(LEAST(a.end_date, w.week_start + 6) -
					 GREATEST(a.start_date, w.week_start) + 1)
				) AS allocated_hours
			FROM selected_people p
			JOIN assignments a ON a.person_id = p.id
			JOIN weeks w
				ON a.start_date <= w.week_start + 6
				AND a.end_date >= w.week_start
			GROUP BY p.id, w.week_start
		)
		SELECT
			p.id,
			p.name,
			p.weekly_hours::double precision,
			array_agg(
				COALESCE(a.allocated_hours, 0)::double precision
				ORDER BY w.week_start
			)
		FROM selected_people p
		CROSS JOIN weeks w
		LEFT JOIN allocations a
			ON a.person_id = p.id AND a.week_start = w.week_start
		GROUP BY p.id, p.name, p.weekly_hours
		ORDER BY p.id
	`, dateRange.From.Format(time.DateOnly), dateRange.To.Format(time.DateOnly), cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("query capacity: %w", err)
	}
	defer rows.Close()

	people := make([]capacity.Person, 0, limit)
	for rows.Next() {
		var person capacity.Person
		if err := rows.Scan(&person.ID, &person.Name, &person.WeeklyHours, &person.AllocatedHours); err != nil {
			return nil, fmt.Errorf("scan capacity: %w", err)
		}
		people = append(people, person)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read capacity rows: %w", err)
	}
	return people, nil
}
