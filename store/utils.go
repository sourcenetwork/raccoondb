package store

import (
	"context"

	"github.com/sourcenetwork/raccoondb/v2/iterator"
)

const batchSize uint = 100

// DeleteAll iterates over kv and deletes all records.
func DeleteAll(ctx context.Context, kv KVStore) error {
	// create it, batchSize keys to memory
	// delete
	// repeat until take / consume yiels no items
	//
	// load
	iter, err := kv.Iterate(ctx, NewOpenIterator())
	if err != nil {
		return err
	}

	for {
		err := iter.Next(ctx)
		if err != nil {
			return err
		}
		if iter.Finished() {
			break
		}

		_, err = kv.Delete(ctx, iter.CurrentKey())
		if err != nil {
			return err
		}
	}

	return nil
}

// IteratePrefix does a prefix Iteration on a Iterable store
// The prefix is added to the opt stary range
// and the result iterator is wrapped in a prefix iterator
// to stop yielding after the prefix is no longer found.
func IteratePrefix[T any](ctx context.Context, iterable Iterable[T], prefix []byte, opt IteratorOpt, stripPrefix bool) (iterator.Iterator[T], error) {
	start := opt.GetLeftBound()
	start = ConcatKey(prefix, start)
	opt = opt.WithLeftBound(start)

	if opt.GetRightBound() != nil {
		end := opt.GetRightBound()
		end = ConcatKey(prefix, end)
		opt.WithRightBound(end)
	}

	iter, err := iterable.Iterate(ctx, opt)
	if err != nil {
		return nil, err
	}
	return iterator.NewPrefixIterator(prefix, iter, stripPrefix), nil
}

// ConcatKey returns a slice which contains the concatination of prefix with key
func ConcatKey(prefix, key []byte) []byte {
	bytes := make([]byte, 0, len(prefix)+len(key))
	bytes = append(bytes, prefix...)
	bytes = append(bytes, key...)
	return bytes
}
