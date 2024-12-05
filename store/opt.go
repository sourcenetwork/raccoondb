package store

// IteratorOpt configures the behavior of an Iterator
// TODO improve this UX
type IteratorOpt struct {
	// start represents the lower bound of iteration
	// If nil will start at the smallest element
	start []byte

	// end represents the uper bound of iteration
	// If nil will end at the largest element
	end []byte

	// reverse iterates the store backwards
	reverse bool
}

// NewOpenIterator returns an Iterator Option which steps through all keys in a store
func NewOpenIterator() IteratorOpt {
	return IteratorOpt{}
}

// NewOpenIterator returns an Iterator Option which steps through all keys in a store
func NewBoundIterator(start, end []byte) IteratorOpt {
	return IteratorOpt{
		start: start,
		end:   end,
	}
}

func (o IteratorOpt) WithReverse(reverse bool) IteratorOpt {
	o.reverse = reverse
	return o
}

func (o IteratorOpt) IsReverse() bool {
	return o.reverse
}

func (o IteratorOpt) WithLeftBound(start []byte) IteratorOpt {
	o.start = start
	return o
}

func (o IteratorOpt) WithRightBound(end []byte) IteratorOpt {
	o.end = end
	return o
}

func (o IteratorOpt) IsOpen() bool {
	return o.start == nil && o.end == nil
}

func (o IteratorOpt) GetRightBound() []byte {
	return o.end
}

func (o IteratorOpt) GetLeftBound() []byte {
	return o.start
}
