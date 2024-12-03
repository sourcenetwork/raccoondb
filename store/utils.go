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
	iter, err := kv.Iterate(ctx, iterator.NewOpenIterator())
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
