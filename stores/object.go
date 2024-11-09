package stores

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/raccoondb/errors"
	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/marshal"
	"github.com/sourcenetwork/raccoondb/types"
)

var ErrKeyObjectStore = errors.New("KeyObjectStore error")

func newErrKeyObject(method string, msg string, err error) error {
	return fmt.Errorf("%w: %v: %v: %w", ErrKeyObjectStore, method, msg, err)
}

// Return a KeyObjectStore from a KVStore using marshaler to (un)marshal objects.
func NewKeyObjectStore[O any](kv KVStore, marshaler marshal.Marshaler[O]) KeyObjectStore[O] {
	countedKV := NewCountedKVStore(kv)
	return KeyObjectStore[O]{
		kv:        countedKV,
		marshaler: marshaler,
	}
}

// KeyObjectStore implements raccoon's ObjKV interface
type KeyObjectStore[Obj any] struct {
	kv        *CountedKVStore
	marshaler marshal.Marshaler[Obj]
}

// Fetch object from store using the given key
func (s *KeyObjectStore[Obj]) Get(ctx context.Context, key []byte) (types.Option[Obj], error) {
	opt, err := s.kv.Get(ctx, key)
	if err != nil {
		return types.None[Obj](), newErrKeyObject("Get", "failed to fetch object", err)
	}
	if opt.Empty() {
		return types.None[Obj](), nil
	}

	obj, err := s.marshaler.Unmarshal(opt.GetValue())
	if err != nil {
		return types.None[Obj](), newErrKeyObject("Get", "failed unmarshaling object", err)
	}

	return types.Some(obj), nil
}

// Set key with obj
func (s *KeyObjectStore[Obj]) Set(ctx context.Context, key []byte, obj Obj) (RecordCreated, error) {
	bytes, err := s.marshaler.Marshal(&obj)
	if err != nil {
		return false, newErrKeyObject("Set", "marshaling object failed", err)
	}
	return s.kv.Set(ctx, key, bytes)
}

// Remove key from store
func (s *KeyObjectStore[Obj]) Delete(ctx context.Context, key []byte) (RecordRemoved, error) {
	removed, err := s.kv.Delete(ctx, key)
	if err != nil {
		return false, newErrKeyObject("Delete", "deleting record", err)
	}
	return removed, nil
}

// Check whether key exists in KVStore
func (s *KeyObjectStore[Obj]) Has(ctx context.Context, key []byte) (bool, error) {
	has, err := s.kv.Has(ctx, key)
	if err != nil {
		return false, newErrKeyObject("Has", "checking", err)
	}
	return has, nil
}

func (s *KeyObjectStore[Obj]) Iterate(ctx context.Context, opts iterator.IteratorOpt) (iterator.Iterator[Obj], error) {
	iter, err := s.kv.Iterate(ctx, opts)
	if err != nil {
		return nil, newErrKeyObject("Iterate", "creating iterator", err)
	}

	objIter := iterator.MapFailable(iter, func(bytes []byte) (Obj, error) {
		return s.marshaler.Unmarshal(bytes)
	})

	return objIter, nil
}

func (s *KeyObjectStore[Obj]) GetObjectCount(ctx context.Context) (uint64, error) {
	count, err := s.kv.GetCount(ctx)
	if err != nil {
		return 0, newErrKeyObject("GetObjectCount", "getting count", err)
	}
	return count, nil
}
