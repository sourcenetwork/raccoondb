package iterator

// IteratorOpt configures the behavior of an Iterator
// TODO improve this UX
type IteratorOpt struct {
	// Start represents the lower bound of iteration
	// If nil will start at the smallest element
	Start []byte

	// End represents the uper bound of iteration
	// If nil will end at the largest element
	End []byte

	// Prefix does a prefix iteration on the store
	Prefix []byte

	// Reverse iterates the store backwards
	Reverse bool
}

func NewOpenIterator() IteratorOpt {
	return IteratorOpt{}
}
