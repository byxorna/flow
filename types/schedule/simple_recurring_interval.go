package schedule

// NOTE(gabe): thanks to https://github.com/victorcoder/dkron/blob/master/cron/constantdelay.go for some bootstrapping

import "time"

// SimpleRecurringInterval is how to express a recurring interval, in simplistic format
// like "every 5 minutes".
// implements Interval
type SimpleRecurringInterval struct {
	Delay time.Duration
}

// Every returns a crontab Schedule that activates once every duration.
// Delays of less than a second are not supported (will round up to 1 second).
// Any fields less than a Second are truncated.
func Every(duration time.Duration) SimpleRecurringInterval {
	if duration < time.Second {
		duration = time.Second
	}
	return SimpleRecurringInterval{
		Delay: duration - time.Duration(duration.Nanoseconds())%time.Second,
	}
}

// Next returns the next time this should be run.
// This rounds so that the next activation time will be on the second.
func (schedule SimpleRecurringInterval) Next(t time.Time) time.Time {
	return t.Add(schedule.Delay - time.Duration(t.Nanosecond())*time.Nanosecond)
}
