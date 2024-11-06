package stores

// TODO
// Add logs
// Wrap Errors

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/types"
)

// Return a KeyObjectStore from a KVStore using marshaler to (un)marshal objects.
func NewKeyObjectStore[O any](kv KVStore, marshaler types.Marshaler[O]) KeyObjectStore[O] {
	countedKV := NewCountedKVStore(kv)
	return KeyObjectStore[O]{
		kv:        countedKV,
		marshaler: marshaler,
	}
}

// KeyObjectStore implements raccoon's ObjKV interface
type KeyObjectStore[Obj any] struct {
	kv        CountedKVStore
	marshaler types.Marshaler[Obj]
}

// Fetch object from store using the given key
func (s *KeyObjectStore[Obj]) Get(ctx context.Context, key []byte) (types.Option[Obj], error) {
	opt, err := s.kv.Get(ctx, key)
	if err != nil {
		err = fmt.Errorf("failed to fetch key %v: %w", key, err)
		return types.None[Obj](), err
	}
	if opt.Empty() {
		return types.None[Obj](), nil
	}

	obj, err := s.marshaler.Unmarshal(opt.GetValue())
	if err != nil {
		err = fmt.Errorf("failed unmarshaling obj from key %v: %w", key, err)
		return types.None[Obj](), err
	}

	return types.Some(obj), nil
}

// Set key with obj
func (s *KeyObjectStore[Obj]) Set(ctx context.Context, key []byte, obj Obj) (RecordCreated, error) {
	bytes, err := s.marshaler.Marshal(&obj)
	if err != nil {
		return false, fmt.Errorf("failed marshaling obj %v: %w", key, err)
	}

	return s.kv.Set(ctx, key, bytes)
}

// Remove key from store
func (s *KeyObjectStore[Obj]) Delete(ctx context.Context, key []byte) (RecordRemoved, error) {
	return s.kv.Delete(ctx, key)
}

// Check whether key exists in KVStore
func (s *KeyObjectStore[Obj]) Has(ctx context.Context, key []byte) (bool, error) {
	return s.kv.Has(ctx, key)
}

func (s *KeyObjectStore[Obj]) Iterate(ctx context.Context, opts iterator.IteratorOpt) (iterator.Iterator[Obj], error) {
	iter, err := s.kv.Iterate(ctx, opts)
	if err != nil {
		return nil, err
	}

	mapper := func(bytes []byte) (Obj, error) {
		return s.marshaler.Unmarshal(bytes)
	}

	objIter := iterator.MapFailable(iter, mapper)
	return objIter, nil
}

func (s *KeyObjectStore[Obj]) GetObjectCount(ctx context.Context) (uint64, error) {
	return s.kv.GetCount(ctx)
}
