package job

type Status uint8

const (
	Pending uint8 = iota
	Running
	Success
	Failed
	PartialyFailed
)

// Status returns the Spec's status of last run
func (j *Spec) Status() Status {
	// TODO get all executions for last run of this job
	// and figure out what status it is...
	return Pending
}
