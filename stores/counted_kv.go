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

var _ KVStore = (*CountedKVStore)(nil)

const countPrefix = "count"
const vals = "objs"
const counterKey = "i"

// Return a KeyObjectStore from a KVStore using marshaler to (un)marshal objects.
func NewCountedKVStore(kv KVStore) CountedKVStore {
	countStore := NewPrefixedKV(kv, []byte(countPrefix))
	counter := NewCounterStore(countStore, nil)
	vals := NewPrefixedKV(kv, []byte(vals))

	return CountedKVStore{
		baseStore: kv,
		vals:      vals,
		counter:   counter,
	}
}

// CountedKVStore implements raccoon's ObjKV interface
type CountedKVStore struct {
	baseStore KVStore
	vals      KVStore
	counter   CounterStore
}

// Fetch object from store using the given key
func (s *CountedKVStore) Get(ctx context.Context, key []byte) (types.Option[[]byte], error) {
	opt, err := s.vals.Get(ctx, key)
	if err != nil {
		err = fmt.Errorf("failed to fetch key %v: %w", key, err)
		return types.None[[]byte](), err
	}

	return opt, nil
}

// Set key with obj
func (s *CountedKVStore) Set(ctx context.Context, key []byte, value []byte) (RecordCreated, error) {
	has, err := s.Has(ctx, key)
	if err != nil {
		return false, err
	}
	if !has {
		err = s.counter.Increment(ctx, []byte(counterKey))
		if err != nil {
			return false, err
		}
	}

	return s.vals.Set(ctx, key, value)
}

// Remove key from store
func (s *CountedKVStore) Delete(ctx context.Context, key []byte) (RecordRemoved, error) {
	has, err := s.Has(ctx, key)
	if err != nil {
		return false, err
	}
	if has {
		err = s.counter.Decrement(ctx, []byte(counterKey))
		if err != nil {
			return false, err
		}
	}

	return s.vals.Delete(ctx, key)
}

// Check whether key exists in KVStore
func (s *CountedKVStore) Has(ctx context.Context, key []byte) (bool, error) {
	return s.vals.Has(ctx, key)
}

func (s *CountedKVStore) Iterate(ctx context.Context, opts iterator.IteratorOpt) (iterator.BytesIterator, error) {
	return s.vals.Iterate(ctx, opts)
}

func (s *CountedKVStore) GetCount(ctx context.Context) (uint64, error) {
	return s.counter.Get(ctx, []byte(counterKey))
}
