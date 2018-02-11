package executor

import (
	"fmt"
)

// Parameters is interface implemented by executors to define the fields a job
// can populate to control how a job is to be run by them. i.e. docker container,
// memory limits, etc
type Parameters interface {
	Type() Type
	Params() map[string]string
}
