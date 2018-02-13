package executor

import (
	"github.com/byxorna/flow/types/execution"
)

// Executor ...
type Executor interface {
	//Run(namespace string, name string) error
	Run(instance *execution.Instance) error
	String() string
	DefaultParameters() (Parameters, error)
}
