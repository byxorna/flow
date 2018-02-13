package shell

import (
	"fmt"
	"time"

	"github.com/byxorna/flow/types/execution"
	"github.com/byxorna/flow/types/executor"
	"github.com/byxorna/flow/types/storage"
	"github.com/sirupsen/logrus"
)

var (
	// ErrWrongExecutor is returned when a job scheduled for another executor is attempted to be run
	ErrWrongExecutor = fmt.Errorf("job is not configured to run with shell executor")
	log              = logrus.WithFields(logrus.Fields{"module": "executor/shell"})
)

// Executor is a shell executor
type Executor struct{}

// Parameters is the type for shell executor parameters
// TODO(gabe) this parameters type feels mad kludgy to me
type Parameters struct {
	// How many concurrent jobs the executor can run
	Concurrency int
}

// New returns a new shell executor
func New(backend *storage.Store) *Executor {
	return &Executor{
		store: backend,
		Settings: Parameters{
			Concurrency: 1,
		},
	}
}

// DefaultParameters is the default params for Shell Executors
func (e *Executor) DefaultParameters() (executor.Parameters, error) {
	return &Parameters{}, nil
}

// Run executes an instance of a job
func (e *Executor) Run(instance *execution.Instance) error {
	//TODO(gabe) make this handle N executions and delegate to workers!
	log.WithFields(logrus.Fields{
		"id":        instance.id,
		"namespace": instance.Namespace,
		"job":       instance.Job,
	}).Infof("executing instance")

	j, err := e.store.GetJob(instance.Namespace, instance.Job)
	if err != nil {
		return err
	}
	if j.Executor != e.Type() {
		return ErrWrongExecutor
	}
	e.StartedAt = time.Now()
	e.store.SetExecution(instance)
	log.Infof("would have run: %s", j.String())

	// when all done, update the execution
	instance.FinishedAt = time.Now()
	//TODO!!!
	instance.Success = true
	instance.Output = []byte{}
	e.store.SetExecution(instance)

	j, _ := e.store.GetJob(instance.Namespace, instance.Job)
	//TODO!!! how do we ensure this locks around a job so we avoid concurrent modification?
	j.SuccessCount += 1
	e.store.SetJob(j)

	return nil
}

// Type ...
func (p *Parameters) Type() executor.Type {
	return executor.TypeShell
}

// Params ...
func (p *Parameters) Params() map[string]string {
	return map[string]string{
		"concurrency": p.Concurrency,
	}
}
