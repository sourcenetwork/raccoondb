package table

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/raccoondb/v2/iterator"
)

// MaterializeObjects receives an iterator of Object Keys
// and returns an iterator which fetches the unmarshaled version of the corresponding objects
func MaterializeObjects[T any](ctx context.Context, store *Table[T], keys ObjKeyIter) iterator.Iterator[T] {
	objIter := iterator.MapFailable(keys, func(key []byte) (T, error) {
		var zero T
		opt, err := store.Get(ctx, key)
		if err != nil {
			return zero, fmt.Errorf("fetching key %v: %w", string(key), err)
		}
		if opt.Empty() {
			return zero, fmt.Errorf("key %v: indexed value was not found in store", string(key))
		}
		return opt.GetValue(), nil
	})
	return objIter
}
