package iterator

import "github.com/sourcenetwork/raccoondb/types"

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
	return nil
}

type filterIterator[T any] struct {
	inner     Iterator[T]
	predicate Predicate[T]
}

func (i *filterIterator[T]) Next() error {
	return i.inner.Next()
}

func (i *filterIterator[T]) Value() types.Option[T] {
	opt := i.inner.Value()
	if opt.Empty() {
		return opt
	}

	val := opt.GetValue()
	if i.predicate(val) {
		return types.None[T]()
	}

	return types.Some(val)
}

func (i *filterIterator[T]) Finished() bool {
	return i.inner.Finished()
}

func (i *filterIterator[T]) Close() error {
	return i.inner.Close()
}

func (i *filterIterator[T]) GetParams() IteratorOpt {
	return i.inner.GetParams()
}

func (i *filterIterator[T]) CurrentKey() []byte {
	return i.CurrentKey()
}

type whileIterator[T any] struct {
	inner     Iterator[T]
	predicate Predicate[T]
	done      bool
}

func (i *whileIterator[T]) Next() error {
	if i.done {
		return nil
	}
	return i.inner.Next()
}

func (i *whileIterator[T]) Value() types.Option[T] {
	opt := i.inner.Value()
	if opt.Empty() {
		return opt
	}

	val := opt.GetValue()
	if !i.predicate(val) {
		i.done = true
		return types.None[T]()
	}

	return types.Some(val)
}

func (i *whileIterator[T]) Finished() bool {
	return i.done || i.inner.Finished()
}

func (i *whileIterator[T]) Close() error {
	return i.inner.Close()
}

func (i *whileIterator[T]) GetParams() IteratorOpt {
	return i.inner.GetParams()
}

func (i *whileIterator[T]) CurrentKey() []byte {
	if i.Finished() {
		return nil
	}
	return i.CurrentKey()
}
