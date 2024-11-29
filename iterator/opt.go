package iterator

// IteratorOpt configures the behavior of an Iterator
// TODO improve this UX
type IteratorOpt struct {
	// start represents the lower bound of iteration
	// If nil will start at the smallest element
	start []byte

	// end represents the uper bound of iteration
	// If nil will end at the largest element
	end []byte

	// prefix does a prefix iteration on the store
	prefix []byte

	// reverse iterates the store backwards
	reverse bool
}

// NewOpenIterator returns an Iterator Option which steps through all keys in a store
func NewOpenIterator() IteratorOpt {
	return IteratorOpt{}
}
