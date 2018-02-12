package executor

import ()

// Executor ...
type Executor interface {
	Run(namespace string, name string) error
	String() string
	DefaultParameters() (Parameters, error)
}
