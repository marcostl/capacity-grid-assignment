package people

import "errors"

var ErrNotFound = errors.New("person not found")

type Person struct {
	ID          int64
	Name        string
	WeeklyHours float64
}
