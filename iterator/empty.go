package iterator

import (
	"context"

	"github.com/sourcenetwork/raccoondb/types"
)

var _ Iterator[any] = (*emptyIterator[any])(nil)

// emptyIterator models an iterator which returns no values
type emptyIterator[T any] struct {
	opt IteratorOpt
}

// NewEmptyIterator returns an iterator which has no values
// ie it's always Finished
func NewEmptyIterator[T any](opt IteratorOpt) Iterator[T] {
	return &emptyIterator[T]{
		opt: opt,
	}
}

func (i *emptyIterator[T]) Next(_ context.Context) error { return nil }
func (i *emptyIterator[T]) Value() types.Option[T]       { return types.None[T]() }
func (i *emptyIterator[T]) CurrentKey() []byte           { return nil }
func (i *emptyIterator[T]) Finished() bool               { return true }
func (i *emptyIterator[T]) Close() error                 { return nil }
func (i *emptyIterator[T]) GetParams() IteratorOpt       { return i.opt }
