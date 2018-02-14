package job

// ID is an identifier for a job
type ID struct {
	// Name is the job name
	Name string `json:"name"`
	// Namespace is the job namespace
	Namespace string `json:"namespace"`
}

func (i *ID) String() string {
	return i.Namespace + "/" + i.Name
}
