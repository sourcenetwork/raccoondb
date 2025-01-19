package iterator

import (
	"bytes"
	"context"

	"github.com/sourcenetwork/raccoondb/v2/types"
)

var _ Iterator[any] = (*PrefixIterator[any])(nil)

// NewPrefixIterator returns a new Iterator which returns only elements which contain prefix.
//
// If stripPrefix is true, the Iterator's CurrentKey()
// will be returned without the given prefix.
//
// During the first Next() call, it seeks the first key which contains prefix.
// Produces items for as long as the key contains prefix or until the iterator finishes
func NewPrefixIterator[T any](prefix []byte, iter Iterator[T], stripPrefix bool) *PrefixIterator[T] {
	return &PrefixIterator[T]{
		prefix:      prefix,
		finished:    false,
		iter:        iter,
		initialized: false,
		stripPrefix: stripPrefix,
	}
}

// PrefixIterator wraps an iterator and steps through it for as long as the key contains
// the given prefix.
type PrefixIterator[T any] struct {
	prefix      []byte
	initialized bool
	finished    bool
	iter        Iterator[T]
	stripPrefix bool
}

func (i *PrefixIterator[T]) Finished() bool {
	return i.finished
}

// Next steps the iterator to the next value
// if the next value does not contain prefix, the scan is done
func (i *PrefixIterator[T]) Next(ctx context.Context) error {
	if i.finished {
		return nil
	}

	if !i.initialized {
		found, err := SeekKeyPrefix(ctx, i.iter, i.prefix)
		if !found {
			i.finished = true
		}
		i.initialized = true

		if err != nil {
			return err
		}
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
	if !i.initialized || i.finished {
		return nil
	}
	if i.stripPrefix {
		return i.iter.CurrentKey()[len(i.prefix):]
	}
	return i.iter.CurrentKey()
}

func (i *PrefixIterator[T]) Value() (types.Option[T], error) {
	if !i.initialized || i.finished {
		return types.None[T](), nil
	}
	return i.iter.Value()
}

func (i *PrefixIterator[T]) Close() error {
	return i.iter.Close()
}
