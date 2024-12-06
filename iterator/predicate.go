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

// While returns an iterator of items for as long as predicte is valid
func While[T any](iter Iterator[T], predicate Predicate[T]) Iterator[T] {
	return &whileIterator[T]{
		inner:     iter,
		predicate: predicate,
		done:      false,
	}
}

type filterIterator[T any] struct {
	inner         Iterator[T]
	predicate     Predicate[T]
	stopOnFailure bool
	value         types.Option[T]
}

func (i *filterIterator[T]) Next(ctx context.Context) error {
	for {
		err := i.inner.Next(ctx)
		if err != nil {
			return err
		}
		opt := i.inner.Value()
		if opt.Empty() {
			return nil
		}
		if i.predicate(opt.GetValue()) {
			return nil
		}
	}
}

func (i *filterIterator[T]) Value() types.Option[T] {
	return i.inner.Value()
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

type whileIterator[T any] struct {
	inner     Iterator[T]
	predicate Predicate[T]
	done      bool
}

func (i *whileIterator[T]) Next(ctx context.Context) error {
	for {
		err := i.inner.Next(ctx)
		if err != nil {
			return err
		}
		opt := i.inner.Value()
		if opt.Empty() {
			return nil
		}
		if i.predicate(opt.GetValue()) {
			return nil
		} else {
			i.done = true
			return nil
		}
	}
}

func (i *whileIterator[T]) Value() types.Option[T] {
	if i.done {
		return types.None[T]()
	}
	return i.inner.Value()
}

func (i *whileIterator[T]) Finished() bool {
	return i.done || i.inner.Finished()
}

func (i *whileIterator[T]) Close() error {
	return i.inner.Close()
}

func (i *whileIterator[T]) CurrentKey() []byte {
	if i.done {
		return nil
	}
	return i.CurrentKey()
}
