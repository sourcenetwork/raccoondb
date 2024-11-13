package composite

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/raccoondb/errors"
	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/stores"
	"github.com/sourcenetwork/raccoondb/types"
)

var ErrIndexedObjectStore = errors.New("indexed object store")

func newIndexedObjectErr(method string, msg string, err error) error {
	return fmt.Errorf("%w: %v: %v: %w", ErrIndexedObjectStore, method, msg, err)
}

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
		return false, newIndexedObjectErr("Delete", "fetching old record", err)
	}
	if opt.Empty() {
		return false, nil
	}
	obj := opt.GetValue()

	err = s.removeFromIndexes(ctx, key, &obj)
	if err != nil {
		return false, newIndexedObjectErr("Delete", "removing from indexes", err)
	}

	_, err = s.objStore.Delete(ctx, key)
	if err != nil {
		return false, newIndexedObjectErr("Delete", "removing record", err)
	}

	return true, nil
}

func (s *IndexedObjectStore[T]) Set(ctx context.Context, key []byte, obj *T) (stores.RecordCreated, error) {
	opt, err := s.objStore.Get(ctx, key)
	if err != nil {
		return false, newIndexedObjectErr("Set", "fetching old record", err)
	}
	if !opt.Empty() {
		s.removeFromIndexes(ctx, key, obj)
	}

	result, err := s.objStore.Set(ctx, key, *obj)
	if err != nil {
		return false, newIndexedObjectErr("Set", "setting record", err)
	}

	for _, idx := range s.indexes {
		i := idx.extractor(obj)
		bucket, err := idx.marshaler.Marshal(&i)
		if err != nil {
			return false, newIndexedObjectErr("Set", "marshaling bucket", err)
		}
		_, err = idx.index.IndexValue(ctx, bucket, key)
		if err != nil {
			return false, newIndexedObjectErr("Set", "setting index", err)
		}
	}

	return result, nil
}

func (s *IndexedObjectStore[T]) removeFromIndexes(ctx context.Context, key []byte, obj *T) error {
	for _, idx := range s.indexes {
		i := idx.extractor(obj)
		bucket, err := idx.marshaler.Marshal(&i)
		if err != nil {
			return fmt.Errorf("marshaling bucket: %w", err)
		}
		_, err = idx.index.RemoveItem(ctx, bucket, key)
		if err != nil {
			return fmt.Errorf("removing from index %v: %w", idx.name, err)
		}
	}
	return nil
}

func (s *IndexedObjectStore[T]) Get(ctx context.Context, key []byte) (types.Option[T], error) {
	opt, err := s.objStore.Get(ctx, key)
	if err != nil {
		return types.None[T](), newIndexedObjectErr("Get", "fetching record", err)
	}
	return opt, nil
}

func (s *IndexedObjectStore[T]) Iterate(ctx context.Context, opt iterator.IteratorOpt) (iterator.Iterator[T], error) {
	iter, err := s.objStore.Iterate(ctx, opt)
	if err != nil {
		return nil, newIndexedObjectErr("Iterate", "creating iterator", err)
	}
	return iter, nil
}

func (s *IndexedObjectStore[T]) Has(ctx context.Context, key []byte) (bool, error) {
	has, err := s.objStore.Has(ctx, key)
	if err != nil {
		return false, newIndexedObjectErr("Has", "checking record", err)
	}
	return has, nil
}

func (s *IndexedObjectStore[T]) MaterializeKeyIter(ctx context.Context, keys iterator.BytesIterator) iterator.Iterator[T] {
	return materializeObjects(ctx, s.objStore, keys)
}
