package shell

import (
	"github.com/byxorna/flow/types/executor"
)

// Parameters ...
type Parameters struct{}

// Type ...
func (p *Parameters) Type() executor.Type {
	return executor.TypeShell
}

// Params ...
func (p *Parameters) Params() map[string]string {
	return map[string]string{}
}
