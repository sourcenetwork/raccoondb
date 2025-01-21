package store

import (
	"context"

	"github.com/sourcenetwork/raccoondb/v2/iterator"
)

const DefaultBatch uint = 100

// DeleteAll lists batchSize keys from the store, deletes all records
// and repeats until all records have been removed
// Note: this operation is not atomic, stops at the first error.
func DeleteAll(ctx context.Context, kv KVStore, batchSize uint) (removed uint, err error) {
	for {
		storeIter, err := kv.Iterate(ctx, NewOpenIterator())
		if err != nil {
			return removed, err
		}
		iter, err := iterator.NewTakeNIterator(storeIter, batchSize)
		if err != nil {
			return removed, err
		}
		keys, err := iterator.ConsumeKeys(ctx, iter)
		if err != nil {
			return removed, err
		}
		if len(keys) == 0 {
			break
		}

		for _, key := range keys {
			_, err := kv.Delete(ctx, key)
			if err != nil {
				return removed, err
			}
			removed++
		}
	}
	return removed, nil
}
