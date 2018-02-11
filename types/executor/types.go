package executor

type Type string

var (
	TypeDefault    Type = "default" // default to shell?
	TypeKubernetes Type = "kubernetes"
	TypeMapReduce  Type = "mapreduce"
	TypeShell      Type = "shell"
	TypeDocker     Type = "docker"
)
