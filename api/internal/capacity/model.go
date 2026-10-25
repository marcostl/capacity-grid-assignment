package capacity

import "time"

const (
	DefaultPageSize = 100
	MaxPageSize     = 250
)

// Request describes a capacity query before its range is normalized to weeks.
type Request struct {
	From   time.Time
	To     time.Time
	Limit  int
	Cursor int64
}

// Range is an inclusive date range.
type Range struct {
	From time.Time
	To   time.Time
}

// Week is an inclusive Monday-through-Sunday range.
type Week struct {
	Start time.Time
	End   time.Time
}

// Person contains one allocation value for each week in Response.Weeks.
type Person struct {
	ID             int64
	Name           string
	WeeklyHours    float64
	AllocatedHours []float64
}

// Response is the application representation of one page of capacity data.
type Response struct {
	Range      Range
	Weeks      []Week
	People     []Person
	Pagination Pagination
}

type Pagination struct {
	Limit      int
	NextCursor *int64
}

// ValidationError is safe for an inbound adapter to show to a user.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
