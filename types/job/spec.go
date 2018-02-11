package job

import (
	"fmt"
	"sync"
	"time"

	"github.com/byxorna/flow/types/executor"
	"github.com/byxorna/flow/types/schedule"
	"github.com/byxorna/flow/types/storage"
	"github.com/docker/libkv/store"
)

// Spec is a Job specification that is provided via API
// to define a job
type Spec struct {
	// Name the name of the job
	Name string
	// Namespace the namespace of the job
	Namespace string

	// Annotations are arbitrary tags associated with the job
	Annotations map[string]string

	// Disabled
	Disabled bool

	// EnvVars are extra env vars to inject into job
	EnvVars map[string]string

	// Owner of the job.
	Owner string `json:"owner"`

	// Owner email of the job.
	OwnerEmail string `json:"owner_email"`

	// Number of successful executions of this job.
	SuccessCount int `json:"success_count"`

	// Number of errors running this job.
	ErrorCount int `json:"error_count"`

	// Last time this job executed succesful.
	LastSuccess time.Time `json:"last_success"`

	// Last time this job failed.
	LastError time.Time `json:"last_error"`

	// Jobs that are dependent upon this one will be run after this job runs.
	DependentJobs []string `json:"dependent_jobs"`

	// Job id of job that this job is dependent upon.
	ParentJob string `json:"parent_job"`

	// Schedule is the desired run schedule
	Schedule schedule.Interval

	// Executor is Which executor to require (if any)
	Executor executor.Type

	// ExecutorParameters are attributes passed to executor to define execution
	// like docker image, entrypoint, memory parameters, etc
	ExecutorParameters executor.Parameters

	// Labels are labels to identify this job
	Labels map[string]string

	// ExecutorConstraints are labels that need satisfied by executor to run this job
	// i.e. specific OS, kernel attributes, etc
	ExecutorConstraints map[string]string

	// storage is the backend storage
	storage storage.Store

	running sync.Mutex
	lock    store.Locker
}

// String is a string rep of a job
func (j *Spec) String() string {
	return fmt.Sprintf("Job %s scheduled at %s with executor %s and labels %v", j.Name, j.Schedule.String(), j.Executor, j.Labels)
}
