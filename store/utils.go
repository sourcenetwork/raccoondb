package store

import (
	"context"

	"github.com/sourcenetwork/raccoondb/v2/iterator"
)

// IteratePrefix does a prefix Iteration on a Iterable store.
//
// IteratePrefix is a more performant option to doing prefix iteration, as opposed to an iterator.PrefixIterator,
// since it can leverage the underlying KV start iteration range.
// In practice it means the start point will be found through binary search, as opposed to a linear scan.
func IteratePrefix[T any](ctx context.Context, iterable Iterable[T], prefix []byte, opt IterationParam, stripPrefix bool) (iterator.Iterator[T], error) {
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
	return iterator.NewPrefixIterator(ctx, prefix, iter, stripPrefix)
}

// ConcatKey returns a slice which contains the concatination of prefix with key
func ConcatKey(prefix, key []byte) []byte {
	bytes := make([]byte, 0, len(prefix)+len(key))
	bytes = append(bytes, prefix...)
	bytes = append(bytes, key...)
	return bytes
}
