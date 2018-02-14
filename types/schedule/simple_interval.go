package schedule

import (
	"fmt"
	"time"
)

// SimpleInterval is a recurring interval that is periodic
type SimpleInterval struct {
	// Period is the period we execute on
	Period time.Duration `json:"period"`
}

// Every duration interval
func Every(duration time.Duration) SimpleInterval {
	if duration < time.Second {
		duration = time.Second
	}
	return SimpleInterval{
		Period: duration - time.Duration(duration.Nanoseconds())%time.Second,
	}
}

// Next returns the next time this should be run.
// This rounds so that the next activation time will be on the second.
func (s SimpleInterval) Next(t time.Time) time.Time {
	return t.Add(s.Period - time.Duration(t.Nanosecond())*time.Nanosecond)
}

// String ...
func (s SimpleInterval) String() string {
	return fmt.Sprintf("every %s", s.Period)
}
