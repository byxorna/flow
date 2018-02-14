package schedule

import (
	"fmt"
	"time"
)

// OneShot is a one-shot job at a give time. non-recurring
type OneShot struct {
	Date time.Time `json:"date"`
}

// At returns a one-shot schedule
func At(date time.Time) OneShot {
	return OneShot{
		Date: date,
	}
}

// Next tells when the job should run next. This is only a one-shot
// though, so if the t is after Date, return the Date.
func (s OneShot) Next(t time.Time) time.Time {
	// if t is in the past, return the date we want to execute at
	if s.Date.After(t) {
		return s.Date
	}
	// else, return a date that will never be reached
	return t.AddDate(69, 0, 0)
}

// String ...
func (s OneShot) String() string {
	return fmt.Sprintf("@%s", s.Date)
}
