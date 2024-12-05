package primitives

// TODO
// Add logs
// Wrap Errors

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/raccoondb/v2/errors"
	"github.com/sourcenetwork/raccoondb/v2/store"
	"github.com/sourcenetwork/raccoondb/v2/types"
)

// ErrCountedKVStore is a top level error for all errors produced by CountedKVStore
var ErrCountedKVStore = errors.New("CountedKVStore error")

var _ store.KVStore = (*countedKVStore)(nil)

const countPrefix = "count/"
const valsPrefix = "vals/"
const counterKey = "i"

// Return an instance of CountedKVStore which tracks the amount of values managed by the store.
func NewCountedKVStore(kv store.KVStore) CountedKVStore {
	vals := NewPrefixedKV(kv, []byte(valsPrefix))
	countPrefixed := NewPrefixedKV(kv, []byte(countPrefix))
	counter := NewCounterStore(countPrefixed)

	return &countedKVStore{
		baseStore: kv,
		vals:      vals,
		counter:   counter,
	}
}

// countedKVStore implements raccoon's CountedKVStore interface
type countedKVStore struct {
	baseStore store.KVStore
	vals      store.KVStore
	counter   CounterStore
}

func (s *countedKVStore) Get(ctx context.Context, key []byte) (types.Option[[]byte], error) {
	opt, err := s.vals.Get(ctx, key)
	if err != nil {
		return types.None[[]byte](), fmt.Errorf("%w: get: %w", ErrCountedKVStore, err)
	}
	return opt, nil
}

func (s *countedKVStore) Set(ctx context.Context, key []byte, value []byte) (store.KeyCreated, error) {
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

func (s *countedKVStore) Delete(ctx context.Context, key []byte) (store.KeyRemoved, error) {
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

func (s *countedKVStore) Has(ctx context.Context, key []byte) (bool, error) {
	has, err := s.vals.Has(ctx, key)
	if err != nil {
		return false, fmt.Errorf("%w: has: %w", ErrCountedKVStore, err)
	}
	return has, nil
}

func (s *countedKVStore) Iterate(ctx context.Context, opts store.IterationParam) (store.StoreIterator[[]byte], error) {
	iter, err := s.vals.Iterate(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("%w: iterate: %w", ErrCountedKVStore, err)
	}
	return iter, nil
}

func (s *countedKVStore) GetCount(ctx context.Context) (uint64, error) {
	i, err := s.counter.Get(ctx, []byte(counterKey))
	if err != nil {
		return 0, fmt.Errorf("%w: getting counter value: %w", ErrCountedKVStore, err)
	}
	return i, nil
}
