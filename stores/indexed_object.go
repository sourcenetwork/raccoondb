package stores

import (
	"context"

	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/types"
)

type ObjectIndex[T any, I any] struct {
}

func (i *ObjectIndex[T, I]) List(ctx context.Context, bucket I) (iterator.Iterator[T], error) {
	return nil, nil
}

func (i *ObjectIndex[T, I]) ListKeys(ctx context.Context, bucket I) (iterator.BytesIterator, error) {
	return nil, nil
}

func (i *ObjectIndex[T, I]) IterateBuckets(ctx context.Context) (iterator.Iterator[I], error) {
	return nil, nil
}

type Index[T any, I any] struct {
	name      string
	marshaler types.Marshaler[I]
	mapper    func(T) I
	store     FieldIndexStore
}

func (i *Index[T, I]) GetObjectIndex(ctx context.Context) *ObjectIndex[T, I] {
	return nil
}

// index decl func to get a readable version of the index store from store

func NewStoreWithIndexes[Obj any](marshaler types.Marshaler[Obj]) IndexedObjectStore[Obj] {
	return IndexedObjectStore[Obj]{}
}

type IndexedObjectStore[T any] struct {
	objStore KeyObjectStore[T]
	indexes  map[string]Index[T, any]
}

func (s *IndexedObjectStore[T]) Initialize(kv KVStore, indexes ...[]*Index[T, any]) error {
	// TODO create substores and so on
	return nil
}

func (s *IndexedObjectStore[T]) Delete(ctx context.Context, key []byte) (RecordRemoved, error) {
	opt, err := s.objStore.Get(ctx, key)
	if err != nil {
		return false, err
	}
	if opt.Empty() {
		return false, nil // do something about noop
	}

	for _, idx := range s.indexes {
		err := idx.store.Delete(ctx, key)
		if err != nil {
			return false, err
		}
	}

	result, err := s.objStore.Delete(ctx, key)
	if err != nil {
		return false, err
	}

	return result, nil
}

func (s *IndexedObjectStore[T]) Set(ctx context.Context, key []byte, obj T) (RecordCreated, error) {
	result, err := s.objStore.Set(ctx, key, obj)
	if err != nil {
		return false, err
	}

	for _, idx := range s.indexes {
		i := idx.mapper(obj)
		bucket, err := idx.marshaler.Marshal(&i)
		if err != nil {
			return false, err
		}
		err = idx.store.IndexValue(ctx, bucket, key)
		if err != nil {
			return false, err
		}
	}
	return result, nil
}

func (s *IndexedObjectStore[T]) UpdateIndexes(ctx context.Context) error {
	return nil

}

func (s *IndexedObjectStore[T]) GetReadStore(ctx context.Context) ReadStore {
	return nil
}
