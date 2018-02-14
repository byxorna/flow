package job

import (
	"fmt"
	"sync"
	"time"

	"github.com/byxorna/flow/types"
	"github.com/byxorna/flow/types/schedule"
)

const (
	// StoragePath ...
	StoragePath = "jobs"
)

var (
	// ErrRequiresSchedule ...
	ErrRequiresSchedule = fmt.Errorf("job requires either a parent job or a schedule")
	// ErrSameParent ...
	ErrSameParent = fmt.Errorf("job cannot be its own parent")
)

// Spec is a Job specification that is provided via API
// to define a job
type Spec struct {
	// ID the name of the job
	ID ID `json:"ID"`

	// Annotations are arbitrary tags associated with the job
	Annotations map[string]string `json:"annotations"`

	// Disabled
	Disabled bool `json:"disabled"`

	// EnvVars are extra env vars to inject into job
	EnvVars map[string]string `json:"env_vars"`

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
	DependentJobs []ID `json:"dependent_jobs,omitempty"`

	// Job id of job that this job is dependent upon.
	ParentJob *ID `json:"parent_job,omitempty"`

	// Schedule is the desired run schedule
	Schedule schedule.Schedule `json:"schedule,omitempty"`

	// Executor is Which executor to require (if any)
	Executor types.Executor `json:"executor"`

	// ExecutorParameters are attributes passed to executor to define execution
	// like docker image, entrypoint, memory parameters, etc
	ExecutorParameters map[string]string `json:"executor_parameters"`

	// ExecutorConstraints are labels that need satisfied by executor to run this job
	// i.e. specific OS, kernel attributes, etc
	ExecutorConstraints map[string]string `json:"constraints"`

	// Labels are labels to identify this job
	Labels map[string]string `json:"labels"`

	running sync.Mutex
}

// String is a string rep of a job
func (j *Spec) String() string {
	return fmt.Sprintf("Job %s scheduled at %s with executor %s and labels %v", j.ID, j.Schedule.String(), j.Executor, j.Labels)
}

// Validate processes a Spec, sets default fields as necessary, and explodes if there is a validation error
func (j *Spec) Validate() error {
	if j.ParentJob != nil && *j.ParentJob == j.ID {
		return ErrSameParent
	}
	if j.ParentJob == nil {
		// require a Schedule
		if j.Schedule == nil {
			return ErrRequiresSchedule
		}
	}

	if j.Executor == types.DefaultExecutor {
		return fmt.Errorf("must specify an explicit executor")
	}
	return nil
}

// Path returns the path to a job given a keyspace
func (j *Spec) Path(keyspace string) string {
	return fmt.Sprintf("%s/%s/%s/%s", keyspace, StoragePath, j.ID.Namespace, j.ID.Name)
}

// Prefix returns the path to a namespaced job in the storage system
func Prefix(keyspace string, id ID) string {
	return fmt.Sprintf("%s/%s/%s/%s", keyspace, StoragePath, id.Namespace, id.Name)
}
