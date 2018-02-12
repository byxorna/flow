package execution

import "fmt"

const (
	// InstancesPath is the path in storage where instances are stored
	InstancesPath = "instances"
)

// Prefix returns the path in the storage layer for a given job's instances
func Prefix(prefix string, ns string, job string) string {
	return fmt.Sprintf("%s/%s/%s/%s", prefix, InstancesPath, ns, job)
}
