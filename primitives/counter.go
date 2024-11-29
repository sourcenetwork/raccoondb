package primitives

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/raccoondb/v2/errors"
	"github.com/sourcenetwork/raccoondb/v2/marshal"
	"github.com/sourcenetwork/raccoondb/v2/store"
)

var ErrCounterStore = errors.New("counter store")

func wrapCounterErr(method string, err error) error {
	return fmt.Errorf("counter %v: %w", method, err)
}

// NewCounterStore returns a new CounterStore
func NewCounterStore(kv store.KVStore) CounterStore {
	return CounterStore{
		kv: kv,
	}
}

// CounterStore abstracts a KVStore to create a store similar to a Program Counter
// where each key has an int value associated to it
type CounterStore struct {
	kv store.KVStore
}

// GetFree returns the next free number in the counter
func (r *CounterStore) GetNext(ctx context.Context, key []byte) (uint64, error) {
	current, err := r.getValue(ctx, key)
	if err != nil {
		return 0, wrapCounterErr("GetNext", err)
	}
	return current + 1, nil
}

// Get return the Counter's current value - 0 if it hasn't been initialized
func (r *CounterStore) Get(ctx context.Context, key []byte) (uint64, error) {
	val, err := r.getValue(ctx, key)
	if err != nil {
		return 0, wrapCounterErr("Get", err)
	}
	return val, nil
}

// Has return true if the counter exists for key
func (r *CounterStore) Has(ctx context.Context, key []byte) (bool, error) {
	has, err := r.kv.Has(ctx, key)
	if err != nil {
		return false, wrapCounterErr("Has", err)
	}
	return has, nil
}

// getValue return the Counter's current value - 0 if it hasn't been initialized
func (r *CounterStore) getValue(ctx context.Context, key []byte) (uint64, error) {
	var currID uint64 = 0
	opt, err := r.kv.Get(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("%w: getting value: %w", ErrCounterStore, err)
	}
	if !opt.Empty() {
		currID = marshal.DecodeUInt(opt.GetValue())
	}
	return currID, nil
}

// Increment updates the counter to the next free number
func (r *CounterStore) setValue(ctx context.Context, key []byte, counter uint64) error {
	_, err := r.kv.Set(ctx, key, marshal.EncodeUInt(counter))
	if err != nil {
		return fmt.Errorf("%w: setting value: %w", ErrCounterStore, err)
	}
	return nil
}

// Increment increments the counter by 1, returns new counter value
func (r *CounterStore) Increment(ctx context.Context, key []byte) (uint64, error) {
	// TODO make it atomic
	free, err := r.GetNext(ctx, key)
	if err != nil {
		return 0, wrapCounterErr("Increment", err)
	}

	err = r.setValue(ctx, key, free)
	if err != nil {
		return 0, wrapCounterErr("Increment", err)
	}

	return free, nil
}

// Decrement reduces the counter by 1, returns new counter value
func (r *CounterStore) Decrement(ctx context.Context, key []byte) (uint64, error) {
	counter, err := r.Get(ctx, key)
	if err != nil {
		return 0, wrapCounterErr("Decrement", err)
	}
	counter -= 1
	err = r.setValue(ctx, key, counter)
	if err != nil {
		return 0, wrapCounterErr("Decrement", err)
	}
	return counter, nil
}

func (r *CounterStore) DeleteCounter(ctx context.Context, key []byte) (store.KeyRemoved, error) {
	removed, err := r.kv.Delete(ctx, key)
	if err != nil {
		return false, wrapCounterErr("DeleteCounter", err)
	}
	return removed, nil
}
