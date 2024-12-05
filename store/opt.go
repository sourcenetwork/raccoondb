package store

// IterationParam sets the params for Iterating over a Store
type IterationParam struct {
	// start represents the lower bound of iteration
	// If nil will start at the smallest element
	start []byte

	// end represents the uper bound of iteration
	// If nil will end at the largest element
	end []byte

	// reverse iterates the store backwards
	reverse bool
}

// NewOpenIterator returns an IterationParam which does a full table scan
func NewOpenIterator() IterationParam {
	return IterationParam{}
}

// NewBoundIterator returns an IteratorParam bound to start and end key
func NewBoundIterator(start, end []byte) IterationParam {
	return IterationParam{
		start: start,
		end:   end,
	}
}

// WithReverse can be set to True in order to do revese iteration
func (o IterationParam) WithReverse(reverse bool) IterationParam {
	o.reverse = reverse
	return o
}

// IsReverse is true if the Iteration is supposed to be reversed
func (o IterationParam) IsReverse() bool {
	return o.reverse
}

// WithLeftBound sets the start point for iteration
func (o IterationParam) WithLeftBound(start []byte) IterationParam {
	o.start = start
	return o
}

// WithRightBound sets the end point for iteration
func (o IterationParam) WithRightBound(end []byte) IterationParam {
	o.end = end
	return o
}

// IsOpen is true if the Iteration is unbound (ie full scan)
func (o IterationParam) IsOpen() bool {
	return o.start == nil && o.end == nil
}

// GetRightBound returns the end interval of the iterator (nil if not set)
func (o IterationParam) GetRightBound() []byte {
	return o.end
}

// GetLeftBound returns the start point of the iterator (nil if not set)
func (o IterationParam) GetLeftBound() []byte {
	return o.start
}
