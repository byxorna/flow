package job

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var (
	// ErrJobDisabled is the error when trying to run a disabled job
	ErrJobDisabled = fmt.Errorf("job is flagged as disabled")
	log            = logrus.WithFields(logrus.Fields{"module": "job"})
)

// Run ...
/*
func (j *Spec) Run() error {
	j.running.Lock()
	defer j.running.Unlock()

	if j.Disabled == false {
		// Check if it's runnable
		if j.isRunnable() {
			log.WithFields(logrus.Fields{
				"namespace": j.Namespace,
				"job":       j.Name,
				"schedule":  j.Schedule.String(),
			}).Debug("enqueuing job")



			// TODO: this should get an existing instance if this is a rerun
			i := execution.NewInstance(j.Namespace, j.Name)
			// TODO: should enqueue jobs to executors instead, and perform job fit?
			exe, err := j.GetExecutor()
			if err != nil {
				return err
			}
			return exe.Run(i)
		}
		return nil
	}
	return ErrJobDisabled
}
*/
