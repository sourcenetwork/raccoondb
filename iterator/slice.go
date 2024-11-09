package iterator

import (
	"github.com/sourcenetwork/raccoondb/marshal"
	"github.com/sourcenetwork/raccoondb/types"
)

func FromSlice[T any](ts []T) Iterator[T] {
	return &SliceAdapter[T]{
		vals: ts,
		idx:  ^uint64(0),
		done: false,
		val:  types.None[T](),
	}
}

var _ Iterator[any] = (*SliceAdapter[any])(nil)

// SliceAdapter wraps a slice which implements the Iterator interface
type SliceAdapter[T any] struct {
	vals []T
	idx  uint64
	val  types.Option[T]
	done bool
}

func (a *SliceAdapter[T]) Next() error {
	if a.done {
		return nil
	}

	a.idx++
	a.val = types.Some(a.vals[a.idx])
	if a.idx == uint64(len(a.vals)-1) {
		a.done = true
	}

	return nil
}

func (a *SliceAdapter[T]) Value() types.Option[T] {
	return a.val
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
	return marshal.EncodeUInt(a.idx)
}
