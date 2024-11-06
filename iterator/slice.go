package iterator

import (
	"github.com/sourcenetwork/raccoondb/marshal"
	"github.com/sourcenetwork/raccoondb/types"
)

func FromSlice[T any](ts []T) Iterator[T] {
	return &SliceAdapter[T]{
		vals: ts,
		idx:  0,
		done: false,
	}
}

var _ Iterator[any] = (*SliceAdapter[any])(nil)

// SliceAdapter wraps a slice which implements the Iterator interface
type SliceAdapter[T any] struct {
	vals []T
	idx  uint64
	done bool
}

func (a *SliceAdapter[T]) Next() error {
	if a.idx+1 == uint64(len(a.vals)) {
		a.done = true
		return nil
	}
	a.idx += 1
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
	return marshal.EncodeUInt(a.idx)
}
