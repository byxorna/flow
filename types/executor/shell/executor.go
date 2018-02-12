package shell

import (
	"github.com/byxorna/flow/types/executor"
)

// Executor is a shell executor
type Executor struct{}

// New returns a new shell executor
func New() *Executor {
	return &Executor{}
}

// Parameters is the params for Shell Executors
type Parameters struct{}

// DefaultParameters is the default params for Shell Executors
func (e *Executor) DefaultParameters() executor.Parameters {
	return &Parameters{}
}

// Type is the executor type this parameters object belongs to
func (p *Parameters) Type() executor.Type {
	return executor.TypeShell
}

// Params are the parameters a shell executor cares about
func (p *Parameters) Params() map[string]string {
	return map[string]string{}
}
