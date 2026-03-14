package grep

type Options struct {
	AContext     int
	BContext     int
	CountSame    bool
	IgnoreCase   bool
	InvertFilter bool
	ExactMatch   bool
	StringNumber bool
}
