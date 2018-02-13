package execution

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Instance is a Job execution instance
type Instance struct {
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

	// ID is a unique ID of this instance
	ID uuid.UUID `json:"id,omitempty"`
}

// NewInstance returns a new execution.Instance
func NewInstance(ns string, name string) *Instance {
	return &Instance{
		Job:       name,
		Namespace: ns,
		Group:     time.Now().UnixNano(),
		Attempt:   1,
		ID:        uuid.New(),
	}
}

// Key is the unique instance of this execution
/*
func (e *Instance) Key() string {
	// TODO does this require some randomness? need a UUID?
	//return fmt.Sprintf("%s-%s-%d", e.Namespace, e.Job, e.StartedAt.UnixNano())
	return e.UUID.String()
}
*/

// String returns a string for this instance
func (e *Instance) String() string {
	return fmt.Sprintf("%s/%s:%s", e.Namespace, e.Job, e.ID.String())
}
