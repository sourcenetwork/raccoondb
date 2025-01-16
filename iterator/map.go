package iterator

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/raccoondb/v2/errors"
	"github.com/sourcenetwork/raccoondb/v2/types"
)

var _ Iterator[any] = (*mapIter[any, any])(nil)

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
	return &mapIter[T, U]{
		inner:  iterator,
		mapper: mapper,
		val:    types.None[U](),
	}
}

// Map applies the Mapper function for every element in the Iterator
func Map[T, U any](iterator Iterator[T], mapper Mapper[T, U]) Iterator[U] {
	m := func(t T) (U, error) {
		return mapper(t), nil
	}
	return &mapIter[T, U]{
		inner:  iterator,
		mapper: m,
		val:    types.None[U](),
	}
}

// mapIter is an iterator which applies a mapping function for every element in the inner iter
type mapIter[T, U any] struct {
	inner       Iterator[T]
	mapper      FailableMapper[T, U]
	mapErr      error
	val         types.Option[U]
	initialized bool
}

func (i *mapIter[T, U]) Next(ctx context.Context) error {
	err := i.inner.Next(ctx)
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

func (i *mapIter[T, U]) Value() types.Option[U] {
	return i.val
}

func (i *mapIter[T, U]) Finished() bool {
	return i.inner.Finished()
}

func (i *mapIter[T, U]) Close() error {
	return i.inner.Close()
}

func (i *mapIter[T, U]) CurrentKey() []byte {
	return i.inner.CurrentKey()
}
