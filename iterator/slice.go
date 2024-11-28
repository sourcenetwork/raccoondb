package iterator

import (
	"github.com/sourcenetwork/raccoondb/marshal"
	"github.com/sourcenetwork/raccoondb/types"
)

var _ Iterator[any] = (*SliceAdapter[any])(nil)

// FromSlice returns wraps a slice with an iterator
// which walks through the slice elems
func FromSlice[T any](ts []T) Iterator[T] {
	return &SliceAdapter[T]{
		vals: ts,
		idx:  0,
		done: len(ts) == 0,
	}
}

// SliceAdapter wraps a slice which implements the Iterator interface
type SliceAdapter[T any] struct {
	vals []T
	idx  uint64
	done bool
}

func (a *SliceAdapter[T]) Next() error {
	if a.done {
		return nil
	}

	a.idx++
	if a.idx == uint64(len(a.vals)) {
		a.done = true
		return nil
	}

	return nil
}

func (a *SliceAdapter[T]) Value() types.Option[T] {
	if a.done {
		return types.None[T]()
	}
	return types.Some(a.vals[a.idx])
}

func (a *SliceAdapter[T]) Finished() bool {
	return a.done
}

func (a *SliceAdapter[T]) Close() error {
	return nil
}

func (a *SliceAdapter[T]) GetParams() IteratorOpt {
	return IteratorOpt{
		Start:   nil,
		End:     nil,
		Prefix:  nil,
		Reverse: false,
	}
}

func (a *SliceAdapter[T]) CurrentKey() []byte {
	if a.done {
		return nil
	}
	return marshal.EncodeUInt(a.idx)
}
