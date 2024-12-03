package table

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/raccoondb/v2/iterator"
	"github.com/sourcenetwork/raccoondb/v2/marshal"
	"github.com/sourcenetwork/raccoondb/v2/primitives"
	"github.com/sourcenetwork/raccoondb/v2/store"
	"github.com/sourcenetwork/raccoondb/v2/types"
)

var _ store.KVStore = (*Table[[]byte])(nil)

func NewTable[T any](kv store.KVStore, marshaler marshal.Marshaler[T]) *Table[T] {
	indexesKv := primitives.NewPrefixedKV(kv, []byte(idxsPrefix))
	objKv := primitives.NewPrefixedKV(kv, []byte(objsPrefix))
	keyObjStore := primitives.NewKeyObjectStore(objKv, marshaler)

	return &Table[T]{
		baseStore: kv,
		objStore:  &keyObjStore,
		indexes:   make(map[string]indexWrite[T]),
		idxsKv:    indexesKv,
	}
}

type Table[T any] struct {
	baseStore store.KVStore
	objStore  *primitives.KeyObjectStore[T]
	indexes   map[string]indexWrite[T]
	idxsKv    store.KVStore
}

func (s *Table[T]) addIndexWriter(name string, writer indexWrite[T]) error {
	_, exists := s.indexes[name]
	if exists {
		return ErrIndexExists
	}
	s.indexes[name] = writer
	return nil
}

func (s *Table[T]) Delete(ctx context.Context, key []byte) (store.KeyRemoved, error) {
	opt, err := s.objStore.Get(ctx, key)
	if err != nil {
		return false, newTableErr("Delete", "fetching old record", err)
	}
	if opt.Empty() {
		return false, nil
	}
	obj := opt.GetValue()

	err = s.removeFromIndexes(ctx, key, &obj)
	if err != nil {
		return false, newTableErr("Delete", "removing from indexes", err)
	}

	_, err = s.objStore.Delete(ctx, key)
	if err != nil {
		return false, newTableErr("Delete", "removing record", err)
	}

	return true, nil
}

func (s *Table[T]) Set(ctx context.Context, key []byte, obj T) (store.KeyCreated, error) {
	opt, err := s.objStore.Get(ctx, key)
	if err != nil {
		return false, newTableErr("Set", "fetching old record", err)
	}
	if !opt.Empty() {
		oldObj := opt.GetValue()
		err := s.removeFromIndexes(ctx, key, &oldObj)
		if err != nil {
			return false, newTableErr("Set", "removing from indexes", err)
		}
	}

	result, err := s.objStore.Set(ctx, key, obj)
	if err != nil {
		return false, newTableErr("Set", "setting record", err)
	}

	for _, idx := range s.indexes {
		_, err := idx.IndexObject(ctx, key, &obj)
		if err != nil {
			return false, newTableErr("Set", "setting index", err)
		}
	}

	return result, nil
}

func (s *Table[T]) removeFromIndexes(ctx context.Context, key []byte, obj *T) error {
	for _, idx := range s.indexes {
		_, err := idx.UnindexObject(ctx, key, obj)
		if err != nil {
			return fmt.Errorf("removing from index %v: %w", idx.GetIndexName(), err)
		}
	}
	return nil
}

func (s *Table[T]) Get(ctx context.Context, key []byte) (types.Option[T], error) {
	opt, err := s.objStore.Get(ctx, key)
	if err != nil {
		return types.None[T](), newTableErr("Get", "fetching record", err)
	}
	return opt, nil
}

func (s *Table[T]) Iterate(ctx context.Context, opt iterator.IteratorOpt) (iterator.Iterator[T], error) {
	iter, err := s.objStore.Iterate(ctx, opt)
	if err != nil {
		return nil, newTableErr("Iterate", "creating iterator", err)
	}
	return iter, nil
}

func (s *Table[T]) Has(ctx context.Context, key []byte) (bool, error) {
	has, err := s.objStore.Has(ctx, key)
	if err != nil {
		return false, newTableErr("Has", "checking record", err)
	}
	return has, nil
}

func (s *Table[T]) MaterializeKeyIter(ctx context.Context, keys ObjKeyIter) iterator.Iterator[T] {
	return MaterializeObjects(ctx, s.objStore, keys)
}

func (s *Table[T]) UpateIndexes(ctx context.Context) error {
	for _, idx := range s.indexes {
		err := idx.Wipe(ctx)
		if err != nil {
			return newTableErr("UpdateIndexes", "wiping indexes", err)
		}
	}

	iter, err := s.objStore.Iterate(ctx, iterator.NewOpenIterator())
	if err != nil {
		return newTableErr("UpdateIndexes", "creating iterator", err)
	}

	for {
		err := iter.Next(ctx)
		if err != nil {
			return newTableErr("UpdateIndexes", "iterating over objects", err)
		}
		if iter.Finished() {
			break
		}
		opt := iter.Value()
		obj := opt.GetValue()
		for _, idx := range s.indexes {
			_, err := idx.IndexObject(ctx, iter.CurrentKey(), &obj)
			if err != nil {
				return newTableErr("UpdateIndexes", "setting index", err)
			}
		}
	}

	return nil
}

func (s *Table[T]) GetCatalogue(ctx context.Context) (*Catalogue, error) {
	dataMap := make(map[string]IndexData)
	for _, idx := range s.indexes {
		buckets, err := idx.GetBucketCount(ctx)
		if err != nil {
			return nil, newTableErr("GetCatalogue", "fetching bucket count", err)
		}
		count, err := idx.GetIndexCount(ctx)
		if err != nil {
			return nil, newTableErr("GetCatalogue", "fetching index count", err)
		}
		data := IndexData{
			Name:               idx.GetIndexName(),
			BucketCount:        buckets,
			IndexedObjectCount: count,
		}
		dataMap[idx.GetIndexName()] = data
	}

	objCount, err := s.objStore.GetObjectCount(ctx)
	if err != nil {
		return nil, newTableErr("GetCatalogue", "fetching obj count", err)
	}

	return &Catalogue{
		ObjectCount: objCount,
		IndexesData: dataMap,
	}, nil
}
