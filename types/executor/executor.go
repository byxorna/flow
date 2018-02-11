package executor

import (
	"github.com/byxorna/flow/types/job"
)

// Parameters is interface implemented by executors to define the fields a job
// can populate to control how a job is to be run by them. i.e. docker container,
// memory limits, etc
type Executor interface {
	Run(j *job.Spec) error
	String() string
	DefaultParameters() (Parameters, error)
}
