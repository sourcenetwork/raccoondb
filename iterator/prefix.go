package iterator

import (
	"bytes"
	"context"
)

var _ Iterator[any] = (*PrefixIterator[any])(nil)

// NewPrefixIterator returns a new Iterator which returns only elements which contain prefix.
//
// If stripPrefix is true, the Iterator's CurrentKey()
// will be returned without the given prefix.
//
// During the first Next() call, it seeks the first key which contains prefix.
// Produces items for as long as the key contains prefix or until the iterator finishes
func NewPrefixIterator[T any](ctx context.Context, prefix []byte, iter Iterator[T], stripPrefix bool) (*PrefixIterator[T], error) {
	prefixIter := &PrefixIterator[T]{
		prefix:      prefix,
		finished:    false,
		iter:        iter,
		stripPrefix: stripPrefix,
	}
	err := prefixIter.initialize(ctx)
	if err != nil {
		return nil, err
	}
	return prefixIter, nil
}

// PrefixIterator wraps an iterator and steps through it for as long as the key contains
// the given prefix.
type PrefixIterator[T any] struct {
	prefix      []byte
	finished    bool
	iter        Iterator[T]
	stripPrefix bool
	zero        T
}

func (i *PrefixIterator[T]) Finished() bool {
	return i.finished
}

func (i *PrefixIterator[T]) initialize(ctx context.Context) error {
	found, err := SeekKeyPrefix(ctx, i.iter, i.prefix)
	if !found {
		i.finished = true
	}

	if err != nil {
		return err
	}
	return nil
}

// Next steps the iterator to the next value
// if the next value does not contain prefix, the scan is done
func (i *PrefixIterator[T]) Next(ctx context.Context) error {
	if i.finished {
		return nil
	}

	err := i.iter.Next(ctx)
	key := i.iter.CurrentKey()
	if i.iter.Finished() || !bytes.HasPrefix(key, i.prefix) {
		i.finished = true
	}

	if err != nil {
		return err
	}

	return nil
}

func (i *PrefixIterator[T]) CurrentKey() []byte {
	if i.finished {
		return nil
	}
	if i.stripPrefix {
		return i.iter.CurrentKey()[len(i.prefix):]
	}
	return i.iter.CurrentKey()
}

func (i *PrefixIterator[T]) Value() (T, error) {
	if i.finished {
		return i.zero, nil
	}
	return i.iter.Value()
}

func (i *PrefixIterator[T]) Close() error {
	return i.iter.Close()
}
