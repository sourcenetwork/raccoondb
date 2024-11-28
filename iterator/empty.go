package iterator

import "github.com/sourcenetwork/raccoondb/types"

var _ Iterator[any] = (*EmptyIterator[any])(nil)

// EmptyIterator models an iterator which returns no values
type EmptyIterator[T any] struct {
	opt IteratorOpt
}

// NewEmptyIterator returns an iterator which has no values
// ie it's always Finished
func NewEmptyIterator[T any](opt IteratorOpt) Iterator[T] {
	return &EmptyIterator[T]{
		opt: opt,
	}
}

func (i *EmptyIterator[T]) Next() error            { return nil }
func (i *EmptyIterator[T]) Value() types.Option[T] { return types.None[T]() }
func (i *EmptyIterator[T]) CurrentKey() []byte     { return nil }
func (i *EmptyIterator[T]) Finished() bool         { return true }
func (i *EmptyIterator[T]) Close() error           { return nil }
func (i *EmptyIterator[T]) GetParams() IteratorOpt { return i.opt }
