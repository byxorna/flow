package executor

type Type string

var (
	TypeKubernetes Type = "kubernetes"
	TypeMapReduce  Type = "mapreduce"
	TypeShell      Type = "shell"
	TypeDocker     Type = "docker"
)
