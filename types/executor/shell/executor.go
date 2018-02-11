package shell

// Executor is a shell executor
type Executor struct {
}

// New returns a new shell executor
func New() *Executor {
	return &Executor{}
}

func (e *Executor) DefaultParameters() Parameters {

}
