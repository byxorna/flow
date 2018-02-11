package execution

import (
	"fmt"
	"time"

	"github.com/byxorna/flow/types/job"
)

type Execution struct {
	// Job name of the job this executions refers to.
	Job string `json:"job_name,omitempty"`
	// Namespace of the job this executions refers to.
	Namespace string `json:"namespace,omitempty"`

	// Start time of the execution.
	StartedAt time.Time `json:"started_at,omitempty"`

	// When the execution finished running.
	FinishedAt time.Time `json:"finished_at,omitempty"`

	// If this execution executed succesfully.
	Success bool `json:"success,omitempty"`

	// Partial output of the execution.
	Output []byte `json:"output,omitempty"`

	// ExecutorAttributes filled by executor (node name, etc)
	ExecutorAttributes map[string]string `json:"executor_attributes,omitempty"`

	// Execution group to what this execution belongs to.
	Group int64 `json:"group,omitempty"`

	// Retry attempt of this execution.
	Attempt uint `json:"attempt,omitempty"`
}

// New returns a new Execution instance
func New(j *job.Spec) *Execution {
	return &Execution{
		Job:       j.Name,
		Namespace: j.Namespace,
		Group:     time.Now().UnixNano(),
		Attempt:   1,
	}
}

// Used to enerate the execution Id
func (e *Execution) Key() string {
	// TODO does this require some randomness?
	return fmt.Sprintf("%s-%s-%d", e.Namespace, e.Job, e.StartedAt.UnixNano())
}
