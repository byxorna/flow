package schedule

// Interval ...
type Interval interface {
	String() string
	Validate() error
}
