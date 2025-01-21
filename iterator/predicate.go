package iterator

import (
	"context"
)

var _ Iterator[any] = (*filterIterator[any])(nil)

type Predicate[T any] func(T) bool

// Filter returns an iterator of items that satisfies the given Predicate
func Filter[T any](iter Iterator[T], predicate Predicate[T]) Iterator[T] {
	return &filterIterator[T]{
		inner:     iter,
		predicate: predicate,
	}
}

type filterIterator[T any] struct {
	inner     Iterator[T]
	predicate Predicate[T]
	value     T
	finished  bool
}

func (i *filterIterator[T]) Next(ctx context.Context) error {
	var zero T
	for {
		if i.inner.Finished() {
			i.finished = true
			i.value = zero
			return nil
		}

		val, err := i.inner.Value()
		if err != nil {
			i.value = zero
			return err
		}

		if i.predicate(val) {
			i.value = val
			return nil
		}
		err = i.inner.Next(ctx)
		if err != nil {
			return err
		}
	}
}

func (i *filterIterator[T]) Value() (T, error) {
	return i.value, nil
}

func (i *filterIterator[T]) Finished() bool {
	return i.inner.Finished()
}

func (i *filterIterator[T]) Close() error {
	return i.inner.Close()
}

func (i *filterIterator[T]) CurrentKey() []byte {
	return i.CurrentKey()
}
