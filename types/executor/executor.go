package executor

import (
	"github.com/byxorna/flow/types/job"
)

// Executor ...
type Executor interface {
	Register(job *job.Spec) error
	Deregister(job *job.Spec) error
	String() string
	DefaultParameters() (Parameters, error)
}
