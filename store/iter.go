package store

import (
	"context"

	"github.com/sourcenetwork/raccoondb/v2/iterator"
	"github.com/sourcenetwork/raccoondb/v2/types"
)

var _ StoreIterator[any] = (*iterAdapter[any])(nil)

// ToStoreIter is an adaptor which returns a StoreIterator from an iterator.Iterator
// The StoreIterator returns the given param when GetParams() is called
func ToStoreIter[T any](iter iterator.Iterator[T], opt IterationParam) StoreIterator[T] {
	return &iterAdapter[T]{
		iter: iter,
		opt:  opt,
	}
}

type iterAdapter[T any] struct {
	iter iterator.Iterator[T]
	opt  IterationParam
}

func (i *iterAdapter[T]) Next(ctx context.Context) error {
	return i.iter.Next(ctx)
}

func (i *iterAdapter[T]) Value() (types.Option[T], error) {
	return i.iter.Value()
}

func (i *iterAdapter[T]) Finished() bool {
	return i.iter.Finished()
}

func (i *iterAdapter[T]) Close() error {
	return i.iter.Close()
}

func (i *iterAdapter[T]) GetParams() IterationParam {
	return i.opt
}

func (i *iterAdapter[T]) CurrentKey() []byte {
	return i.iter.CurrentKey()
}
