package primitives

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/raccoondb/v2/errors"
	"github.com/sourcenetwork/raccoondb/v2/iterator"
	"github.com/sourcenetwork/raccoondb/v2/store"
)

// ErrFieldIndex is a top level error for all errors produced by FieldIndexStore
var ErrFieldIndex = errors.New("FieldIndexStore")

const bucketsPrefix = "buckets/"
const idxPrefix = "idx/"
const bucketCounterPrefix = "bucket_counter/"

func newFieldIndexErr(method string, msg string, err error) error {
	return fmt.Errorf("%w: %v: %v: %w", ErrFieldIndex, method, msg, err)
}

// NewFieldIndexStore returns a new FieldIndexStore from a KVStore
func NewFieldIndexStore(kv store.KVStore) *FieldIndexStore {
	return &FieldIndexStore{
		baseKv:        kv,
		idx:           NewCountedKVStore(NewPrefixedKV(kv, []byte(idxPrefix))),
		buckets:       NewCountedKVStore(NewPrefixedKV(kv, []byte(bucketsPrefix))),
		bucketCounter: NewCounterStore(NewPrefixedKV(kv, []byte(bucketCounterPrefix))),
	}
}

// FieldIndexStore models an indexing structure which groups a set of items under a bucket
// It provides no abstraction over the byte sequences, only a way to iterate over items in a bucket,
// as well as counters for buckets
type FieldIndexStore struct {
	baseKv store.KVStore
	// buckets store the user defined buckets
	buckets CountedKVStore
	// bucketCounter stores a count of elements per bucket
	bucketCounter CounterStore
	// idx stores the indexed values inside each bucket
	idx CountedKVStore
}

// IndexValue adds a value to a bucket
//
// If value was previously inserted in a different bucket, it doesn't scan the index
// to remove it, that is the callers responsability.
func (s *FieldIndexStore) IndexValue(ctx context.Context, bucket []byte, item []byte) (store.KeyCreated, error) {
	_, err := s.buckets.Set(ctx, bucket, bucket)
	if err != nil {
		return false, newFieldIndexErr("IndexValue", "creating bucket", err)
	}

	key := getIdxKey(bucket, item)
	created, err := s.idx.Set(ctx, key, item)
	if err != nil {
		return false, newFieldIndexErr("IndexValue", "indexing value", err)
	}
	if created {
		_, err := s.bucketCounter.Increment(ctx, bucket)
		if err != nil {
			return false, newFieldIndexErr("IndexValue", "incrementing bucket counter", err)
		}
	}

	return created, nil
}

// Has returns true if the given bucket contains item
func (s *FieldIndexStore) Has(ctx context.Context, bucket, item []byte) (bool, error) {
	idxKey := getIdxKey(bucket, item)
	has, err := s.idx.Has(ctx, idxKey)
	if err != nil {
		return false, newFieldIndexErr("Has", "fetching", err)
	}
	return has, nil
}

// IterateBucketValues returns an iterator which returns all values contained in a bucket
func (s *FieldIndexStore) IterateBucketItems(ctx context.Context, bucket []byte, opt store.IterationParam) (iterator.BytesIterator, error) {
	bucketStore := NewPrefixedKV(s.idx, concatKey(bucket, []byte("/")))
	iter, err := bucketStore.Iterate(ctx, opt)
	if err != nil {
		return nil, newFieldIndexErr("IterateBucketItems", "creating iterator", err)
	}
	return iter, nil
}

// RemoveItem removes the given item from bucket
// If bucket did not contain item, return RecordRemoved false
// If the removed item was the last item from the bucket, removes the bucket
func (s *FieldIndexStore) RemoveItem(ctx context.Context, bucket, item []byte) (store.KeyRemoved, error) {
	key := getIdxKey(bucket, item)
	removed, err := s.idx.Delete(ctx, key)
	if err != nil {
		return false, newFieldIndexErr("RemoveItem", "removing item", err)
	}

	removeBucket := false
	if removed {
		count, err := s.bucketCounter.Decrement(ctx, bucket)
		if err != nil {
			return false, newFieldIndexErr("RemoveItem", "decrementing bucket counter", err)
		}
		if count == 0 {
			removeBucket = true
		}
	}

	if removeBucket {
		_, err := s.bucketCounter.Delete(ctx, bucket)
		if err != nil {
			return false, newFieldIndexErr("RemoveItem", "removing bucket counter", err)
		}

		_, err = s.buckets.Delete(ctx, bucket)
		if err != nil {
			return false, newFieldIndexErr("RemoveItem", "removing bucket", err)
		}
	}

	return removed, nil
}

// IterateBuckets returns an iterator over all buckets which currently contains at least
// one item, where the iterator key and its value are the bucket
func (s *FieldIndexStore) IterateBuckets(ctx context.Context, opt store.IterationParam) (iterator.BytesIterator, error) {
	iter, err := s.buckets.Iterate(ctx, opt)
	if err != nil {
		return nil, newFieldIndexErr("IterateBuckets", "creating iterator", err)
	}
	return iter, nil
}

// GetBucketCount returns the number of active buckets
func (s *FieldIndexStore) GetBucketCount(ctx context.Context) (uint64, error) {
	count, err := s.buckets.GetCount(ctx)
	if err != nil {
		return 0, newFieldIndexErr("GetBucketCount", "getting count", err)
	}
	return count, nil
}

// GEtCountItemsInBucket returns the number of items contained in a bucket
func (s *FieldIndexStore) GetCountItemsInBucket(ctx context.Context, bucket []byte) (uint64, error) {
	count, err := s.bucketCounter.Get(ctx, bucket)
	if err != nil {
		return 0, newFieldIndexErr("GetCountItemsInBucket", "getting count", err)
	}
	return count, nil
}

// GetIndexedItemsCount returns the total number of items that have been indexed
// accross all buckets
func (s *FieldIndexStore) GetIndexedItemsCount(ctx context.Context) (uint64, error) {
	count, err := s.idx.GetCount(ctx)
	if err != nil {
		return 0, newFieldIndexErr("GetIndexedItemsCount", "getting count", err)
	}
	return count, nil
}

// Wipe removes all entries from FieldIndexStore
func (s *FieldIndexStore) Wipe(ctx context.Context) error {
	err := store.DeleteAll(ctx, s.baseKv)
	if err != nil {
		return newFieldIndexErr("Wipe", "DeleteAll", err)
	}
	return nil
}

func getIdxKey(bucket []byte, value []byte) []byte {
	l := len(bucket) + len(value) + 1
	key := make([]byte, 0, l)
	key = append(key, bucket...)
	key = append(key, byte('/'))
	key = append(key, value...)
	return key
}
