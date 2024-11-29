package table

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/raccoondb/errors"
	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/marshal"
	"github.com/sourcenetwork/raccoondb/primitives"
	"github.com/sourcenetwork/raccoondb/store"
	"github.com/sourcenetwork/raccoondb/types"
)

var ErrIndexedObjectStore = errors.New("indexed object store")
var ErrIndexExists = errors.New("index already defined")

const objsPrefix = "objs/"
const idxsPrefix = "indexes/"

func NewTable[T any](kv store.KVStore, marshaler marshal.Marshaler[T]) *Table[T] {
	indexesKv := primitives.NewPrefixedKV(kv, []byte(idxsPrefix))
	objKv := primitives.NewPrefixedKV(kv, []byte(objsPrefix))
	keyObjStore := primitives.NewKeyObjectStore(objKv, marshaler)

	return &Table[T]{
		objStore: &keyObjStore,
		indexes:  make(map[string]IndexWriter[T]),
		idxsKv:   indexesKv,
	}
}

func newTableErr(method string, msg string, err error) error {
	return fmt.Errorf("%w: %v: %v: %w", ErrIndexedObjectStore, method, msg, err)
}

type Table[T any] struct {
	objStore *primitives.KeyObjectStore[T]
	indexes  map[string]IndexWriter[T]
	idxsKv   store.KVStore
}

func (s *Table[T]) addIndexWriter(name string, writer IndexWriter[T]) error {
	_, exists := s.indexes[name]
	if exists {
		return ErrIndexExists
	}
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

func (s *Table[T]) Set(ctx context.Context, key []byte, obj *T) (store.KeyCreated, error) {
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
		_, err := idx.IndexObject(ctx, key, obj)
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

func (s *Table[T]) MaterializeKeyIter(ctx context.Context, keys iterator.BytesIterator) iterator.Iterator[T] {
	return MaterializeObjects(ctx, s.objStore, keys)
}

func (s *Table[T]) UpateIndexes(ctx context.Context) error {
	panic("TODO")
}

func (s *Table[T]) GetCatalogue(ctx context.Context) (Catalogue, error) {
	panic("TODO")
}
