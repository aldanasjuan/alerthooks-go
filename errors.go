package alerthooks

import "errors"

var (
	ErrMissingID        = errors.New("params.ID must not be empty")
	ErrInvalidURL       = errors.New("url has to be an absolute url, no localhost allowed, https only and the domain must exist")
	ErrInvalidDueDate   = errors.New("the record due date must be in the future")
	ErrInvalidMethod    = errors.New("the record method must be one of: GET, POST, PUT, PATCH, DELETE")
	ErrInvalidType      = errors.New("the record type must be one of: one_time, persistent, recurring")
	ErrNilRecurring     = errors.New("recurring record should have a recurring object with minutes, hours, days and months")
	ErrMissingRecurring = errors.New("recurring record should have at least one option of minutes, hours, days and months")
	ErrRecurringType    = errors.New("recurring type should be 'week' or 'month'")
	ErrBadMinutes       = errors.New("recurring minutes should have at least one option and should be between 0-59")
	ErrBadHours         = errors.New("recurring hours should have at least one option and should be between 0-23")
	ErrBadDays          = errors.New("recurring days (week or month) should have at least one option and should be between 0-6(weekdays) or 1-31(monthdays)")
	ErrBadMonths        = errors.New("recurring minutes should have at least one option and should be between 1-12")
)
