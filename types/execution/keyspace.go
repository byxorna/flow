package execution

import (
	"fmt"
	"github.com/byxorna/flow/types/job"
)

const (
	// InstancesPath is the path in storage where instances are stored
	InstancesPath = "instances"
)

// Path returns the path in the storage layer for a given job's instances
//func Path(prefix string, ns string, job string) string {
func Path(prefix string, id job.ID) string {
	//return fmt.Sprintf("%s/%s/%s/%s", prefix, InstancesPath, ns, job)
	return fmt.Sprintf("%s/%s/%s/%s", prefix, InstancesPath, id.Namespace, id.Name)
}
