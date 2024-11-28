package table

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/raccoondb/errors"
	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/marshal"
	"github.com/sourcenetwork/raccoondb/primitives"
	"github.com/sourcenetwork/raccoondb/store"
)

type IndexValueExtractor[T, I any] func(*T) I

var ErrObjectIndex = errors.New("ObjectIndex")

func NewIndex[T, I any](table *Table[T], name string, extractor IndexValueExtractor[T, I], marshaler marshal.Marshaler[I]) (IndexReader[T, I], error) {
	idxKv := primitives.NewPrefixedKV(table.idxsKv, []byte(name))
	fieldIdx := primitives.NewFieldIndexStore(idxKv)
	idx := &TableIndex[T, I]{
		name:      name,
		index:     &fieldIdx,
		extractor: extractor,
		marshaler: marshaler,
	}
	err := table.addIndexWriter(name, idx)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

func newIdxErr(method string, msg string, err error) error {
	return fmt.Errorf("%w: %v: %v: %w", ErrObjectIndex, method, msg, err)
}

func newObjectIndexStore[T, I any](name string, idx *primitives.FieldIndexStore, extractor IndexValueExtractor[T, I], marshaler marshal.Marshaler[I]) *TableIndex[T, I] {
	return &TableIndex[T, I]{
		name:      name,
		index:     idx,
		extractor: extractor,
		marshaler: marshaler,
	}
}

type IndexWriter[T any] interface {
	GetIndexName() string
	IndexObject(ctx context.Context, key []byte, obj *T) (store.KeyCreated, error)
	UnindexObject(ctx context.Context, key []byte, obj *T) (store.KeyRemoved, error)
	UpdateIndex(ctx context.Context, key []byte, old *T, new *T) error
}

type IndexReader[T, I any] interface {
	GetIndexName() string
	IterateKeys(ctx context.Context, bucket *I) (ObjKeyIter, error)
	IterateBuckets(ctx context.Context) (iterator.Iterator[I], error)
	GetBucketCount(ctx context.Context) (uint64, error)
	GetIndexCount(ctx context.Context) (uint64, error)
}

var _ IndexWriter[any] = (*TableIndex[any, any])(nil)

type TableIndex[T any, I any] struct {
	name      string
	index     *primitives.FieldIndexStore
	extractor IndexValueExtractor[T, I]
	marshaler marshal.Marshaler[I]
}

// Iterate returns an iterator which steps though the object keys indexed in a given bucket / value
func (i *TableIndex[T, I]) IterateKeys(ctx context.Context, bucket *I) (ObjKeyIter, error) {
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
func (i *TableIndex[T, I]) IterateBuckets(ctx context.Context) (iterator.Iterator[I], error) {
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
func (i *TableIndex[T, I]) GetBucketCount(ctx context.Context) (uint64, error) {
	count, err := i.index.GetBucketCount(ctx)
	if err != nil {
		return 0, newIdxErr("GetBucketCount", "getting count", err)
	}
	return count, nil
}

// GetIndexCount returns the count of the number of objects that have been indexed
func (i *TableIndex[T, I]) GetIndexCount(ctx context.Context) (uint64, error) {
	count, err := i.index.GetIndexedItemsCount(ctx)
	if err != nil {
		return 0, newIdxErr("GetIndexCount", "getting count", err)
	}
	return count, nil
}

// IndexObject indexes the given object using the given key as record id
func (i *TableIndex[T, I]) IndexObject(ctx context.Context, key []byte, obj *T) (store.KeyCreated, error) {
	val := i.extractor(obj)
	bucket, err := i.marshaler.Marshal(&val)
	if err != nil {
		return false, newTableErr("IndexObject", "marshaling bucket", err)
	}

	created, err := i.index.IndexValue(ctx, bucket, key)
	if err != nil {
		return false, newTableErr("IndexObject", "setting index", err)
	}
	return created, nil
}

func (i *TableIndex[T, I]) UpdateIndex(ctx context.Context, key []byte, old *T, new *T) error {
	val := i.extractor(old)
	bucket, err := i.marshaler.Marshal(&val)
	if err != nil {
		return newTableErr("UpdateIndex", "marshaling bucket", err)
	}

	_, err = i.index.RemoveItem(ctx, bucket, key)
	if err != nil {
		return newTableErr("UpdateIndex", "removing old from index", err)
	}

	_, err = i.IndexObject(ctx, key, new)
	if err != nil {
		return newTableErr("UpdateIndex", "indexing new", err)
	}
	return nil
}

func (i *TableIndex[T, I]) UnindexObject(ctx context.Context, key []byte, obj *T) (store.KeyRemoved, error) {
	val := i.extractor(obj)
	bucket, err := i.marshaler.Marshal(&val)
	if err != nil {
		return false, newTableErr("UnindexObject", "marshaling bucket", err)
	}

	removed, err := i.index.RemoveItem(ctx, bucket, key)
	if err != nil {
		return false, newTableErr("UnindexObject", "removing from index", err)
	}
	return removed, nil
}

func (i *TableIndex[T, I]) GetIndexName() string {
	return i.name
}
