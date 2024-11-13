package composite

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/raccoondb/errors"
	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/marshal"
	"github.com/sourcenetwork/raccoondb/stores"
)

type IndexValueExtractor[T, I any] func(*T) I

var ErrObjectIndex = errors.New("ObjectIndex")

func newIdxErr(method string, msg string, err error) error {
	return fmt.Errorf("%w: %v: %v: %w", ErrObjectIndex, method, msg, err)
}

func newObjectIndexStore[T, I any](name string, idx *stores.FieldIndexStore, extractor IndexValueExtractor[T, I], marshaler marshal.Marshaler[I]) *ObjectIndexStore[T, I] {
	return &ObjectIndexStore[T, I]{
		name:      name,
		index:     idx,
		extractor: extractor,
		marshaler: marshaler,
	}
}

type ObjectIndexStore[T any, I any] struct {
	name      string
	index     *stores.FieldIndexStore
	extractor IndexValueExtractor[T, I]
	marshaler marshal.Marshaler[I]
}

// Iterate returns an iterator which steps though the object keys indexed in a given bucket / value
func (i *ObjectIndexStore[T, I]) IterateKeys(ctx context.Context, bucket *I) (ObjKeyIter, error) {
	bytes, err := i.marshaler.Marshal(bucket)
	if err != nil {
		return nil, newIdxErr("IterateKeys", "marshaling bucket", err)
	}

	iter, err := i.index.IterateBucketItems(ctx, bytes)
	if err != nil {
		return nil, newIdxErr("IterateKeys", "buckets iterator", err)
	}
	return iter, nil
}

// IterateValues returns an iterator which steps though the buckets / values in the index
func (i *ObjectIndexStore[T, I]) IterateValues(ctx context.Context) (iterator.Iterator[I], error) {
	iter, err := i.index.IterateBuckets(ctx)
	if err != nil {
		return nil, newIdxErr("IterateValues", "creating iterator", err)
	}

	valIter := iterator.MapFailable(iter, func(bytes []byte) (I, error) {
		var zero I
		val, err := i.marshaler.Unmarshal(bytes)
		if err != nil {
			return zero, err
		}
		return val, nil
	})
	return valIter, nil
}

// GetBucketCount returns the number of buckets / values the index contains
func (i *ObjectIndexStore[T, I]) GetBucketCount(ctx context.Context) (uint64, error) {
	count, err := i.index.GetBucketCount(ctx)
	if err != nil {
		return 0, newIdxErr("GetBucketCount", "getting count", err)
	}
	return count, nil
}

// GetIndexCount returns the count of the number of objects that have been indexed
func (i *ObjectIndexStore[T, I]) GetIndexCount(ctx context.Context) (uint64, error) {
	count, err := i.index.GetIndexedItemsCount(ctx)
	if err != nil {
		return 0, newIdxErr("GetIndexCount", "getting count", err)
	}
	return count, nil
}
