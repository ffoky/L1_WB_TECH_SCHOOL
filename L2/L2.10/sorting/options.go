package sorting

// Options параметры сортировки GNU sort
type Options struct {
	Column            int
	TrimTailSpace     bool
	Reverse           bool
	Numeric           bool
	Unique            bool
	Month             bool
	IsSorted          bool
	HumanReadableSize bool
}
