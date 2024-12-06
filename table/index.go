package table

import (
	"context"

	"github.com/sourcenetwork/raccoondb/v2/iterator"
	"github.com/sourcenetwork/raccoondb/v2/marshal"
	"github.com/sourcenetwork/raccoondb/v2/primitives"
	"github.com/sourcenetwork/raccoondb/v2/store"
)

var _ indexWrite[any] = (*tableIndex[any, any])(nil)
var _ IndexReader[any, any] = (*tableIndex[any, any])(nil)

// IndexValueExtract models a function to be called over an object of type T
// to yield an index value of type I
type IndexValueExtractor[T, I any] func(*T) I

// NewIndex creates a new Index and attaches it to table.
// The created index is updated with every mutation to table.
// Returns a queryable view of the created Index
//
// New indexes aren't automatically updated by the Table.
func NewIndex[T, I any](table *Table[T], name string, extractor IndexValueExtractor[T, I], marshaler marshal.Marshaler[I]) (IndexReader[T, I], error) {
	idxKv := primitives.NewPrefixedKV(table.idxsKv, []byte(name))
	fieldIdx := primitives.NewFieldIndexStore(idxKv)
	idx := &tableIndex[T, I]{
		name:      name,
		index:     fieldIdx,
		extractor: extractor,
		marshaler: marshaler,
	}
	err := table.addIndexWriter(name, idx)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

// indexWriter models an interface to mutate indexes with any bucket type
type indexWrite[T any] interface {
	IndexCatalogue

	// IndexObject indexes obj with the given key in a bucket extrated from obj
	IndexObject(ctx context.Context, key []byte, obj *T) (store.KeyCreated, error)

	// UnindexObject removes obj with the given key from the bucket
	// extracted from it
	UnindexObject(ctx context.Context, key []byte, obj *T) (store.KeyRemoved, error)

	// UpdatedIndex removes obj with the given key from the bucket
	// extracted from old and adds it to new
	UpdateIndex(ctx context.Context, key []byte, old *T, new *T) error

	// Wipe removes all entries from the index
	Wipe(ctx context.Context) error
}

type IndexCatalogue interface {
	// GetIndexName returns the index name
	GetIndexName() string

	// GetBucketCount returns the number of buckets which contain
	// at least one object
	GetBucketCount(ctx context.Context) (uint64, error)

	// GetIndexCount returns the total number of indexed objects
	GetIndexCount(ctx context.Context) (uint64, error)
}

// IndexReader models a readable interface of an index attached to a table
type IndexReader[T, I any] interface {
	IndexCatalogue

	// IterateKeys returns an iterator which yiels the
	// keys of all objects indexed under bucket
	IterateKeys(ctx context.Context, bucket *I, opt store.IterationParam) (ObjKeyIter, error)

	// IterateBuckets returns an iterator which yields all buckets
	// which contains at least one object
	IterateBuckets(ctx context.Context, opt store.IterationParam) (iterator.Iterator[I], error)

	// Has returns true if bucket contains the given key
	Has(ctx context.Context, bucket *I, key []byte) (bool, error)

	// Iterate walks through all elements in index
	Iterate(ctx context.Context, opt store.IterationParam) (store.StoreIterator[[]byte], error)
}

// tableIndex wraps FieldIndexStore abstracting the step of marshaling
// bucket values into bytes and extracting buckes from objects
type tableIndex[T any, I any] struct {
	name      string
	index     *primitives.FieldIndexStore
	extractor IndexValueExtractor[T, I]
	marshaler marshal.Marshaler[I]
}

// Iterate returns an iterator which steps though the object keys indexed in a given bucket / value
func (i *tableIndex[T, I]) IterateKeys(ctx context.Context, bucket *I, opt store.IterationParam) (ObjKeyIter, error) {
	bytes, err := i.marshaler.Marshal(bucket)
	if err != nil {
		return nil, newIndexErr("IterateKeys", "marshaling bucket", err)
	}

	iter, err := i.index.IterateBucketItems(ctx, bytes, opt)
	if err != nil {
		return nil, newIndexErr("IterateKeys", "buckets iterator", err)
	}
	return iter, nil
}

// IterateValues returns an iterator which steps though the buckets / values in the index
func (i *tableIndex[T, I]) IterateBuckets(ctx context.Context, opt store.IterationParam) (iterator.Iterator[I], error) {
	iter, err := i.index.IterateBuckets(ctx, opt)
	if err != nil {
		return nil, newIndexErr("IterateValues", "creating iterator", err)
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
func (i *tableIndex[T, I]) GetBucketCount(ctx context.Context) (uint64, error) {
	count, err := i.index.GetBucketCount(ctx)
	if err != nil {
		return 0, newIndexErr("GetBucketCount", "getting count", err)
	}
	return count, nil
}

// GetIndexCount returns the count of the number of objects that have been indexed
func (i *tableIndex[T, I]) GetIndexCount(ctx context.Context) (uint64, error) {
	count, err := i.index.GetIndexedItemsCount(ctx)
	if err != nil {
		return 0, newIndexErr("GetIndexCount", "getting count", err)
	}
	return count, nil
}

// IndexObject indexes the given object using the given key as record id
func (i *tableIndex[T, I]) IndexObject(ctx context.Context, key []byte, obj *T) (store.KeyCreated, error) {
	val := i.extractor(obj)
	bucket, err := i.marshaler.Marshal(&val)
	if err != nil {
		return false, newIndexErr("IndexObject", "marshaling bucket", err)
	}

	created, err := i.index.IndexValue(ctx, bucket, key)
	if err != nil {
		return false, newIndexErr("IndexObject", "setting index", err)
	}
	return created, nil
}

// UpdateIndex removes key from the index value / bucket of old and moves it to new
func (i *tableIndex[T, I]) UpdateIndex(ctx context.Context, key []byte, old *T, new *T) error {
	val := i.extractor(old)
	bucket, err := i.marshaler.Marshal(&val)
	if err != nil {
		return newIndexErr("UpdateIndex", "marshaling bucket", err)
	}

	_, err = i.index.RemoveItem(ctx, bucket, key)
	if err != nil {
		return newIndexErr("UpdateIndex", "removing old from index", err)
	}

	_, err = i.IndexObject(ctx, key, new)
	if err != nil {
		return newIndexErr("UpdateIndex", "indexing new", err)
	}
	return nil
}

// UnindexObject removes key from the bucket associated to obj
func (i *tableIndex[T, I]) UnindexObject(ctx context.Context, key []byte, obj *T) (store.KeyRemoved, error) {
	val := i.extractor(obj)
	bucket, err := i.marshaler.Marshal(&val)
	if err != nil {
		return false, newIndexErr("UnindexObject", "marshaling bucket", err)
	}

	removed, err := i.index.RemoveItem(ctx, bucket, key)
	if err != nil {
		return false, newIndexErr("UnindexObject", "removing from index", err)
	}
	return removed, nil
}

// GetIndexName returns the index name
func (i *tableIndex[T, I]) GetIndexName() string {
	return i.name
}

// Has returns true if the bucket cointains the given object
func (i *tableIndex[T, I]) Has(ctx context.Context, bucket *I, key []byte) (bool, error) {
	bucketBytes, err := i.marshaler.Marshal(bucket)
	if err != nil {
		return false, newIndexErr("Has", "marshaling bucket", err)
	}

	has, err := i.index.Has(ctx, bucketBytes, key)
	if err != nil {
		return false, newIndexErr("Has", "setting index", err)
	}
	return has, nil
}

func (i *tableIndex[T, I]) Wipe(ctx context.Context) error {
	err := i.index.Wipe(ctx)
	if err != nil {
		return newIndexErr("Wipe", "wipe all", err)
	}
	return nil
}

func (i *tableIndex[T, I]) Iterate(ctx context.Context, opt store.IterationParam) (store.StoreIterator[[]byte], error) {
	iter, err := i.index.Iterate(ctx, opt)
	if err != nil {
		return nil, newIndexErr("Iterate", "", err)
	}
	return iter, nil
}
