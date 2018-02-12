package executor

// Type is the type of executor
type Type string

var (
	// TypeKubernetes ...
	TypeKubernetes Type = "kubernetes"
	// TypeMapReduce ...
	TypeMapReduce Type = "mapreduce"
	// TypeShell ...
	TypeShell Type = "shell"
	// TypeDocker ...
	TypeDocker Type = "docker"
)
