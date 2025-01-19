package iterator

import (
	"context"

	"github.com/sourcenetwork/raccoondb/v2/types"
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
	value     types.Option[T]
	finished  bool
}

func (i *filterIterator[T]) Next(ctx context.Context) error {
	for {
		if i.inner.Finished() {
			i.finished = true
			i.value = types.None[T]()
			return nil
		}

		opt, err := i.inner.Value()
		if err != nil {
			i.value = types.None[T]()
			return err
		}

		if !opt.Empty() {
			// if iterator has no value
			// skip since we are interested on things that match
			// the predicate
			continue
		}

		if i.predicate(opt.GetValue()) {
			i.value = types.Some(opt.GetValue())
			return nil
		}
		err = i.inner.Next(ctx)
		if err != nil {
			return err
		}
	}
}

func (i *filterIterator[T]) Value() (types.Option[T], error) {
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
