package iterator

import (
	"context"
)

var _ Iterator[any] = (*filterIterator[any])(nil)

type Predicate[T any] func(T) bool

// Filter returns an iterator of items that satisfies the given Predicate.
func Filter[T any](iter Iterator[T], predicate Predicate[T]) Iterator[T] {
	f := &filterIterator[T]{
		inner:     iter,
		predicate: predicate,
	}
	// Initialize: position on first matching item
	_ = f.seekNextMatch(context.Background())
	return f
}

type filterIterator[T any] struct {
	inner     Iterator[T]
	predicate Predicate[T]
	value     T
	zero      T
}

// seekNextMatch advances through the inner iterator until it finds
// an item that matches the predicate, or the iterator is exhausted.
func (i *filterIterator[T]) seekNextMatch(ctx context.Context) error {
	for {
		if i.inner.Finished() {
			i.value = i.zero
			return nil
		}

		val, err := i.inner.Value()
		if err != nil {
			i.value = i.zero
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

func (i *filterIterator[T]) Next(ctx context.Context) error {
	if i.inner.Finished() {
		return nil
	}

	// Advance past current match
	err := i.inner.Next(ctx)
	if err != nil {
		return err
	}

	// Find next matching item
	return i.seekNextMatch(ctx)
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
	if i.inner.Finished() {
		return nil
	}
	return i.inner.CurrentKey()
}
