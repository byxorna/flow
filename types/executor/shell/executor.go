package shell

//TODO remove this whole executor when we have something better?
// this is really only useful for debugging

import (
	"fmt"
	"sync"
	"time"

	"github.com/byxorna/flow/types/execution"
	"github.com/byxorna/flow/types/executor"
	"github.com/byxorna/flow/types/job"
	"github.com/byxorna/flow/types/storage"
	"github.com/sirupsen/logrus"
)

var (
	// ErrWrongExecutor is returned when a job scheduled for another executor is attempted to be run
	ErrWrongExecutor = fmt.Errorf("job is not configured to run with shell executor")
	// ErrJobNotFound ...
	ErrJobNotFound = fmt.Errorf("job not found in queue")
	log            = logrus.WithFields(logrus.Fields{"module": "executor/shell"})
)

// Executor is a shell executor
type Executor struct {
	sync.Mutex
	queue []*job.Spec
	store *storage.Store
}

// Parameters is the type for shell executor parameters
// TODO(gabe) this parameters type feels mad kludgy to me
type Parameters struct {
	// How many concurrent jobs the executor can run
	Concurrency int
}

// New returns a new shell executor
func New(backend *storage.Store) (*Executor, error) {

	log.Debug("loading all jobs")
	jobs, err := backend.GetJobs()
	if err != nil {
		return nil, err
	}

	queue := []*job.Spec{}
	for _, j := range jobs {
		if j.Executor == executor.TypeShell {
			queue = append(queue, j)
		}
	}

	log.WithFields(logrus.Fields{"jobs": len(queue)}).Debug("loaded jobs")

	return &Executor{
		store: backend,
		queue: queue,
		Settings: Parameters{
			Concurrency: 1,
		},
	}
}

// Deregister deregisters a job from the queue
func (e *Executor) Deregister(j *spec.Job) error {
	e.Lock()
	defer e.Unlock()
	found := false
	for i, x := range e.queue {
		if x.ID() == j.ID() {
			// splice x out of the queue
			e.queue = append(e.queue[:i], e.queue[i+1:]...)
			return nil
		}
	}
	return ErrJobNotFound
}

// Register registers a job to be processed
func (e *Executor) Register(j *spec.Job) error {
	e.Lock()
	defer e.Unlock()
	if j.Executor == executor.TypeShell {
		e.queue = append(e.queue, j)
	} else {
		return ErrWrongExecutor
	}
}

// DefaultParameters is the default params for Shell Executors
func (e *Executor) DefaultParameters() (executor.Parameters, error) {
	return &Parameters{}, nil
}

/*
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
*/

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
