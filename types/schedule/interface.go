package schedule

import (
	"time"
)

// Schedule is a schedule to execute a job on
type Schedule interface {
	Next(time.Time) time.Time
	String() string
}
