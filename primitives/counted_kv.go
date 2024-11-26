package primitives

// TODO
// Add logs
// Wrap Errors

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/raccoondb/errors"
	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/store"
	"github.com/sourcenetwork/raccoondb/types"
)

var ErrCountedKVStore = errors.New("CounterKVStore error")

var _ store.KVStore = (*CountedKVStore)(nil)

const countPrefix = "count"
const valsPrefix = "vals"
const counterKey = "i"

// Return a KeyObjectStore from a store.KVStore using marshaler to (un)marshal objects.
func NewCountedKVStore(kv store.KVStore) *CountedKVStore {
	countStore := NewPrefixedKV(kv, []byte(countPrefix))
	counter := NewCounterStore(countStore)
	vals := NewPrefixedKV(kv, []byte(valsPrefix))

	return &CountedKVStore{
		baseStore: kv,
		vals:      vals,
		counter:   counter,
	}
}

// CountedKVStore implements raccoon's ObjKV interface
type CountedKVStore struct {
	baseStore store.KVStore
	vals      store.KVStore
	counter   CounterStore
}

// Fetch object from store using the given key
func (s *CountedKVStore) Get(ctx context.Context, key []byte) (types.Option[[]byte], error) {
	opt, err := s.vals.Get(ctx, key)
	if err != nil {
		return types.None[[]byte](), fmt.Errorf("%w: get: %w", ErrCountedKVStore, err)
	}
	return opt, nil
}

// Set key with obj
func (s *CountedKVStore) Set(ctx context.Context, key []byte, value []byte) (store.RecordCreated, error) {
	created, err := s.vals.Set(ctx, key, value)
	if err != nil {
		return false, fmt.Errorf("%w: set: setting record: %w", ErrCountedKVStore, err)
	}
	if created {
		_, err := s.counter.Increment(ctx, []byte(counterKey))
		if err != nil {
			return false, fmt.Errorf("%w: set: incrementing record counter: %w", ErrCountedKVStore, err)
		}
	}
	return created, nil
}

// Remove key from store
func (s *CountedKVStore) Delete(ctx context.Context, key []byte) (store.RecordRemoved, error) {
	removed, err := s.vals.Delete(ctx, key)
	if err != nil {
		return false, fmt.Errorf("%w: delete: deleting record: %w", ErrCountedKVStore, err)
	}
	if removed {
		_, err := s.counter.Decrement(ctx, []byte(counterKey))
		if err != nil {
			return false, fmt.Errorf("%w: delete: decrementing record counter: %w", ErrCountedKVStore, err)
		}
	}
	return removed, nil
}

// Check whether key exists in store.KVStore
func (s *CountedKVStore) Has(ctx context.Context, key []byte) (bool, error) {
	has, err := s.vals.Has(ctx, key)
	if err != nil {
		return false, fmt.Errorf("%w: has: %w", ErrCountedKVStore, err)
	}
	return has, nil
}

func (s *CountedKVStore) Iterate(ctx context.Context, opts iterator.IteratorOpt) (iterator.BytesIterator, error) {
	iter, err := s.vals.Iterate(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("%w: iterate: %w", ErrCountedKVStore, err)
	}
	return iter, nil
}

func (s *CountedKVStore) GetCount(ctx context.Context) (uint64, error) {
	i, err := s.counter.Get(ctx, []byte(counterKey))
	if err != nil {
		return 0, fmt.Errorf("%w: getting counter value: %w", ErrCountedKVStore, err)
	}
	return i, nil
}
