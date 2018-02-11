package shell

import (
	"github.com/byxorna/flow/types/executor"
)

type Parameters struct{}

func (p *Parameters) Type() executor.Type {
	return executor.TypeShell
}

func (p *Parameters) Params() map[string]string {
	return map[string]string{}
}
