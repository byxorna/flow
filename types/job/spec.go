package job

import (
	"fmt"
	"sync"
	"time"

	"github.com/byxorna/flow/types/executor"
	"github.com/byxorna/flow/types/schedule"
	//	"github.com/byxorna/flow/types/storage"
)

var (
	ErrRequiresSchedule = fmt.Errorf("Job requires either a parent job or a schedule!")
)

// Spec is a Job specification that is provided via API
// to define a job
type Spec struct {
	// Name the name of the job
	Name string `json:"name"`
	// Namespace the namespace of the job
	Namespace string `json:"namespace"`

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
	DependentJobs []string `json:"dependent_jobs"`

	// Job id of job that this job is dependent upon.
	ParentJob string `json:"parent_job"`

	// Schedule is the desired run schedule
	Schedule schedule.Interval `json:"schedule"`

	// Executor is Which executor to require (if any)
	Executor executor.Type `json:"executor"`

	// ExecutorParameters are attributes passed to executor to define execution
	// like docker image, entrypoint, memory parameters, etc
	ExecutorParameters executor.Parameters `json:"executor_parameters"`

	// Labels are labels to identify this job
	Labels map[string]string `json:"labels"`

	// ExecutorConstraints are labels that need satisfied by executor to run this job
	// i.e. specific OS, kernel attributes, etc
	ExecutorConstraints map[string]string `json:"constraints"`

	// storage is the backend storage
	//storage storage.Store

	running sync.Mutex
}

// String is a string rep of a job
func (j *Spec) String() string {
	return fmt.Sprintf("Job %s scheduled at %s with executor %s and labels %v", j.Name, j.Schedule.String(), j.Executor, j.Labels)
}

// Validate processes a Spec, sets default fields as necessary, and explodes if there is a validation error
func (j *Spec) Validate() error {
	if j.ParentJob == "" {
		// require a Schedule
		if j.Schedule == nil {
			return ErrRequiresSchedule
		}
		if err := j.Schedule.Validate(); err != nil {
			return err
		}
	}

	if j.Executor == nil {
		j.Executor = executor.TypeShell
	}
	if j.ExecutorParameters == nil {
		j.ExecutorParameters = j.Executor.DefaultParameters()
	}
}
