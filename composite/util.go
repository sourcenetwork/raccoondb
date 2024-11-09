package composite

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/stores"
)

func materializeObjects[T any](ctx context.Context, store *stores.KeyObjectStore[T], keys ObjKeyIter) iterator.Iterator[T] {
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
