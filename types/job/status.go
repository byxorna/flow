package job

// Status ...
type Status uint8

const (
	// Pending ...
	Pending uint8 = iota
	// Running ...
	Running
	// Success ...
	Success
	// Failed ...
	Failed
	// PartiallyFailed ...
	PartiallyFailed
)

// Status returns the Spec's status of last run
func (j *Spec) Status() Status {
	// TODO get all executions for last run of this job
	// and figure out what status it is...
	return Pending
}

func (j *Spec) isRunnable() bool {
	switch status := j.Status(); status {
	case Running:
		return false
	default:
		return true
	}
}
