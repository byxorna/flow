package execution

import (
	"fmt"
	"time"

	"github.com/byxorna/flow/types/job"
	"github.com/google/uuid"
)

// Instance is a Job execution instance
type Instance struct {
	// ID of job
	Job job.ID `json:"job"`
	/*
		// Job name of the job this executions refers to.
		Job string `json:"job_name,omitempty"`
		// Namespace of the job this executions refers to.
		Namespace string `json:"namespace,omitempty"`
	*/

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

	// ID is a unique ID of this instance
	ID uuid.UUID `json:"id,omitempty"`
}

// NewInstance returns a new execution.Instance
func NewInstance(id job.ID) *Instance {
	return &Instance{
		Job:     id,
		Group:   time.Now().UnixNano(),
		Attempt: 1,
		ID:      uuid.New(),
	}
}

// String returns a string for this instance
func (e *Instance) String() string {
	return fmt.Sprintf("%s/%s:%s", e.Job.Namespace, e.Job.Name, e.ID.String())
}
