package table

import "github.com/sourcenetwork/raccoondb/v2/types"

// BucketIterationParam sets the params for Iterating over a Store
type BucketIterationParam[T any] struct {
	// start represents the lower bound of iteration
	// If nil will start at the smallest element
	start types.Option[T]

	// end represents the uper bound of iteration
	// If nil will end at the largest element
	end types.Option[T]

	// reverse iterates the store backwards
	reverse bool
}

// NewOpenIterator returns an IterationParam which does a full table scan
func NewOpenIterator[T any]() BucketIterationParam[T] {
	return BucketIterationParam[T]{
		start:   types.None[T](),
		end:     types.None[T](),
		reverse: false,
	}
}

// NewBoundIterator returns an IteratorParam bound to start and end key
func NewBoundIterator[T any](start, end T) BucketIterationParam[T] {
	param := NewOpenIterator[T]()
	return param.WithLeftBound(start).WithRightBound(end)
}

// WithReverse can be set to True in order to do revese iteration
func (o BucketIterationParam[T]) WithReverse(reverse bool) BucketIterationParam[T] {
	o.reverse = reverse
	return o
}

// IsReverse is true if the Iteration is supposed to be reversed
func (o BucketIterationParam[T]) IsReverse() bool {
	return o.reverse
}

// WithLeftBound sets the start point for iteration
func (o BucketIterationParam[T]) WithLeftBound(start T) BucketIterationParam[T] {
	o.start = types.Some(start)
	return o
}

// WithRightBound sets the end point for iteration
func (o BucketIterationParam[T]) WithRightBound(end T) BucketIterationParam[T] {
	o.end = types.Some(end)
	return o
}

// IsOpen is true if the Iteration is unbound (ie full scan)
func (o BucketIterationParam[T]) IsOpen() bool {
	return o.start.Empty() && o.end.Empty()
}

// GetRightBound returns the end interval of the iterator (none if not set)
func (o BucketIterationParam[T]) GetRightBound() types.Option[T] {
	return o.end
}

// GetLeftBound returns the start point of the iterator (none if not set)
func (o BucketIterationParam[T]) GetLeftBound() types.Option[T] {
	return o.start
}
