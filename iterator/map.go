package iterator

import (
	"fmt"

	"github.com/sourcenetwork/raccoondb/errors"
	"github.com/sourcenetwork/raccoondb/types"
)

// ErrMapper models an error that happened when MapIter attempted to map an element
var ErrMapper = errors.New("mapping failed")

// FailableMapper maps T to U or errors
type FailableMapper[T, U any] func(T) (U, error)

// FailableMapper maps T to U
type Mapper[T, U any] func(T) U

// MapFailable wraps an iterator, applying the mapper function for each
// value of the inner iterator.
// If the mapping for an element fails Next() will return an error and
// Value() will return None.
// Then the Iter will move to the next element
func MapFailable[T, U any](iterator Iterator[T], mapper FailableMapper[T, U]) Iterator[U] {
	return &MapIter[T, U]{
		inner:  iterator,
		mapper: mapper,
	}
}

// Map applies the Mapper function for every element in the Iterator
func Map[T, U any](iterator Iterator[T], mapper Mapper[T, U]) Iterator[U] {
	m := func(t T) (U, error) {
		return mapper(t), nil
	}
	return &MapIter[T, U]{
		inner:  iterator,
		mapper: m,
	}
}

var _ Iterator[any] = (*MapIter[any, any])(nil)

type MapIter[T, U any] struct {
	inner  Iterator[T]
	mapper FailableMapper[T, U]
	mapErr error
	val    types.Option[U]
}

func (i *MapIter[T, U]) Next() error {
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
		key := i.inner.CurrentKey()
		i.val = types.None[U]()
		return fmt.Errorf("elem '%v': %w: %w", string(key), ErrMapper, err)
	}

	i.val = types.Some(u)
	return nil
}

func (i *MapIter[T, U]) Value() types.Option[U] {
	return i.val
}

func (i *MapIter[T, U]) Finished() bool {
	return i.inner.Finished()
}

func (i *MapIter[T, U]) Close() error {
	return i.Close()
}

func (i *MapIter[T, U]) GetParams() IteratorOpt {
	return i.GetParams()
}

func (i *MapIter[T, U]) CurrentKey() []byte {
	return i.CurrentKey()
}
