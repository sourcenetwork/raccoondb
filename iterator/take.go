package iterator

import (
	"context"

	"github.com/sourcenetwork/raccoondb/v2/errors"
)

var _ Iterator[any] = (*TakeNIter[any])(nil)

var ErrZeroCount error = errors.New("take n iter: count must be greater than 0")

// TakeNIter wraps an iterator yielding a maximum number of pre-set values
type TakeNIter[T any] struct {
	iter     Iterator[T]
	n        uint
	max      uint
	finished bool
	zero     T
}

func NewTakeNIterator[T any](iter Iterator[T], count uint) (Iterator[T], error) {
	if count == 0 {
		return nil, ErrZeroCount
	}
	return &TakeNIter[T]{
		iter:     iter,
		max:      count,
		finished: false,
		n:        1,
	}, nil
}

func (i *TakeNIter[T]) Finished() bool {
	return i.finished
}

// Next steps the iterator to the next value
// if the next value does not contain prefix, the scan is done
func (i *TakeNIter[T]) Next(ctx context.Context) error {
	i.n++
	if i.n > i.max {
		i.finished = true
		return nil
	}

	err := i.iter.Next(ctx)
	if i.iter.Finished() {
		i.finished = true
	}

	if err != nil {
		return err
	}

	return nil
}

func (i *TakeNIter[T]) CurrentKey() []byte {
	if i.finished {
		return nil
	}
	return i.iter.CurrentKey()
}

func (i *TakeNIter[T]) Value() (T, error) {
	if i.finished {
		return i.zero, nil
	}
	return i.iter.Value()
}

func (i *TakeNIter[T]) Close() error {
	return i.iter.Close()
}
