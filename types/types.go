package types

// Executor is the type of executor
type Executor string

var (
	// DefaultExecutor is the default executor type
	DefaultExecutor Executor
	// KubernetesExecutor ...
	KubernetesExecutor Executor = "kubernetes"
	// MapReduceExecutor ...
	MapReduceExecutor Executor = "mapreduce"
	// ShellExecutor ...
	ShellExecutor Executor = "shell"
	// MesosExecutor ...
	MesosExecutor Executor = "mesos"
)
