package iterator

import (
	"bytes"

	"github.com/sourcenetwork/raccoondb/types"
)

var _ Iterator[any] = (*PrefixIterator[any])(nil)

func NewPrefixIterator[T any](prefix []byte, iter Iterator[T]) *PrefixIterator[T] {
	firstKey := iter.CurrentKey()
	// if they first key doesn't have the prefix, iterator should be empty and we are done.
	// it means there are no values in the store that satisfy the prefix, that is the precondition
	if !bytes.HasPrefix(firstKey, prefix) {
		return &PrefixIterator[T]{
			finished: true,
			iter:     iter,
			prefix:   prefix,
		}
	}
	return &PrefixIterator[T]{
		prefix:   prefix,
		finished: false,
		iter:     iter,
	}
}

type PrefixIterator[T any] struct {
	prefix   []byte
	finished bool
	iter     Iterator[T]
}

func (i *PrefixIterator[T]) Finished() bool {
	return i.finished
}

// Next steps the iterator to the next value
// if the next value does not contain prefix, the scan is done
func (i *PrefixIterator[T]) Next() error {
	err := i.iter.Next()
	key := i.iter.CurrentKey()
	if i.iter.Finished() || !bytes.HasPrefix(key, i.prefix) {
		i.finished = true
	}

	if err != nil {
		return err
	}

	return nil
}

func (i *PrefixIterator[T]) CurrentKey() (key []byte) {
	if i.finished {
		return nil
	}
	return i.iter.CurrentKey()
}

func (i *PrefixIterator[T]) Value() types.Option[T] {
	if i.finished {
		return types.None[T]()
	}
	return i.iter.Value()
}

func (i *PrefixIterator[T]) Close() error {
	return i.iter.Close()
}

func (i *PrefixIterator[T]) GetParams() IteratorOpt {
	return i.iter.GetParams()
}
