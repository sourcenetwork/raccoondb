package composite

import (
	"context"

	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/stores"
)

func newIndexedObjectStore[T any](store *stores.KeyObjectStore[T], indexes []*ObjectIndexStore[T, any]) IndexedObjectStore[T] {
	return IndexedObjectStore[T]{
		objStore: store,
		indexes:  indexes,
	}
}

type IndexedObjectStore[T any] struct {
	objStore *stores.KeyObjectStore[T]
	indexes  []*ObjectIndexStore[T, any]
}

func (s *IndexedObjectStore[T]) Delete(ctx context.Context, key []byte) (stores.RecordRemoved, error) {
	opt, err := s.objStore.Get(ctx, key)
	if err != nil {
		return false, err
	}
	if opt.Empty() {
		return false, nil // do something about noop
	}
	obj := opt.GetValue()

	err = s.removeFromIndexes(ctx, key, &obj)
	if err != nil {
		return false, err
	}

	_, err = s.objStore.Delete(ctx, key)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *IndexedObjectStore[T]) Set(ctx context.Context, key []byte, obj *T) (stores.RecordCreated, error) {
	opt, err := s.objStore.Get(ctx, key)
	if err != nil {
		return false, err
	}
	if !opt.Empty() {
		s.removeFromIndexes(ctx, key, obj)
	}

	result, err := s.objStore.Set(ctx, key, *obj)
	if err != nil {
		return false, err
	}

	for _, idx := range s.indexes {
		i := idx.extractor(obj)
		bucket, err := idx.marshaler.Marshal(&i)
		if err != nil {
			return false, err
		}
		_, err = idx.index.IndexValue(ctx, bucket, key)
		if err != nil {
			return false, err
		}
	}

	return result, nil
}

func (s *IndexedObjectStore[T]) removeFromIndexes(ctx context.Context, key []byte, obj *T) error {
	for _, idx := range s.indexes {
		i := idx.extractor(obj)
		bucket, err := idx.marshaler.Marshal(&i)
		if err != nil {
			return err
		}
		_, err = idx.index.RemoveItem(ctx, bucket, key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *IndexedObjectStore[T]) UpdateIndexes(ctx context.Context) error {
	return nil

}

func (s *IndexedObjectStore[T]) GetReadStore(ctx context.Context) stores.ReadStore {
	return nil
}

func (s *IndexedObjectStore[T]) MaterializeKeyIter(ctx context.Context, keys iterator.BytesIterator) iterator.Iterator[T] {
	return materializeObjects(ctx, s.objStore, keys)
}
