package iterator

import (
	"fmt"

	"github.com/sourcenetwork/raccoondb/types"
)

type Mapper[T, U any] func(T) U

// Map wraps an iterator, applying the mapper function for each
// value of the inner iterator
func Map[T, U any](iterator Iterator[T], mapper Mapper[T, U]) Iterator[U] {
	return &mapIterator[T, U]{
		inner:  iterator,
		mapper: mapper,
	}
}

// TryMap wraps an iterator, applying the mapper function for each
// value of the inner iterator. If the mapping for an element fails,
// Next() will return an error, Value() will return None and that element is skipped.
func TryMap[T, U any](iterator Iterator[T], mapper FailableMapper[T, U]) Iterator[U] {
	return nil
}

var _ Iterator[any] = (*mapIterator[any, any])(nil)

type mapIterator[T, U any] struct {
	inner  Iterator[T]
	mapper Mapper[T, U]
}

func (i *mapIterator[T, U]) Next() error {
	return i.inner.Next()
}

func (i *mapIterator[T, U]) Value() types.Option[U] {
	opt := i.inner.Value()
	if opt.Empty() {
		return types.None[U]()
	}

	val := opt.GetValue()
	u := i.mapper(val)
	return types.Some(u)
}

func (i *mapIterator[T, U]) Finished() bool {
	return i.inner.Finished()
}

func (i *mapIterator[T, U]) Close() error {
	return i.Close()
}

func (i *mapIterator[T, U]) GetParams() IteratorOpt {
	return i.GetParams()
}

func (i *mapIterator[T, U]) CurrentKey() []byte {
	return i.CurrentKey()
}

var _ Iterator[any] = (*tryMapIter[any, any])(nil)

type tryMapIter[T, U any] struct {
	inner  Iterator[T]
	mapper FailableMapper[T, U]
	mapErr error
	val    types.Option[U]
}

func (i *tryMapIter[T, U]) Next() error {
	err := i.inner.Next()
	if err != nil {
		i.val = types.None[U]()
		return err
	}

	opt := i.inner.Value()
	if opt.Empty() {
		i.val = types.None[U]()
		return nil
	}

	val := opt.GetValue()
	u, err := i.mapper(val)
	if err != nil {
		i.val = types.None[U]()
		return fmt.Errorf("mapping elem %v: %w", err)
	}

	i.val = types.Some(u)
	return nil
}

func (i *tryMapIter[T, U]) Value() types.Option[U] {
	return i.val
}

func (i *tryMapIter[T, U]) Finished() bool {
	return i.inner.Finished()
}

func (i *tryMapIter[T, U]) Close() error {
	return i.Close()
}

func (i *tryMapIter[T, U]) GetParams() IteratorOpt {
	return i.GetParams()
}

func (i *tryMapIter[T, U]) CurrentKey() []byte {
	return i.CurrentKey()
}
