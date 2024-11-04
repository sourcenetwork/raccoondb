package stores

import (
	"context"
	"encoding/binary"

	"github.com/sourcenetwork/raccoondb/types"
)

// NewCounterStore returns a new counter store
func NewCounterStore(kv KVStore, logger types.Logger) CounterStore {
	return CounterStore{
		kv:     kv,
		logger: logger,
	}
}

// CounterStore wraps a KV Store and creates a monotomically increasing counter
type CounterStore struct {
	kv     KVStore
	logger types.Logger
}

func (r *CounterStore) getStore() KVStore {
	return r.kv
}

// GetFree returns the next free number in the counter
func (r *CounterStore) GetNext(ctx context.Context, key []byte) (uint64, error) {
	current, err := r.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	return current + 1, nil
}

// Get return the Counter's current value - 0 if it hasn't been initialized
func (r *CounterStore) Get(ctx context.Context, key []byte) (uint64, error) {
	return r.getValue(ctx, key)
}

// getValue return the Counter's current value - 0 if it hasn't been initialized
func (r *CounterStore) getValue(ctx context.Context, key []byte) (uint64, error) {
	kv := r.getStore()

	var currID uint64 = 0
	opt, err := kv.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	if !opt.Empty() {
		currID = decode(opt.GetValue())
	}
	return currID, nil
}

// Increment updates the counter to the next free number
func (r *CounterStore) setValue(ctx context.Context, key []byte, counter uint64) error {
	kv := r.getStore()

	_, err := kv.Set(ctx, key, encode(counter))
	if err != nil {
		return err
	}

	return nil
}

// Increment increments the counter by 1
func (r *CounterStore) Increment(ctx context.Context, key []byte) error {
	_, err := r.GetNextAndIncrement(ctx, key)
	return err
}

// Decrement reduces the counter by 1
func (r *CounterStore) Decrement(ctx context.Context, key []byte) error {
	counter, err := r.Get(ctx, key)
	if err != nil {
		return err
	}
	counter -= 1
	return r.setValue(ctx, key, counter)
}

// GetNextAndIncrement gets the next free counter and increments it
func (r *CounterStore) GetNextAndIncrement(ctx context.Context, key []byte) (uint64, error) {
	// TODO make it atomic
	free, err := r.GetNext(ctx, key)
	if err != nil {
		return 0, err
	}

	err = r.setValue(ctx, key, free)
	if err != nil {
		return 0, err
	}

	return free, nil
}

// encode maps a uin64 to a big endian byte slice
func encode(counter uint64) []byte {
	buff := make([]byte, 8)
	binary.BigEndian.PutUint64(buff, counter)
	return buff
}

// decode converts a big endian byte slice into a uint64
func decode(counter []byte) uint64 {
	return binary.BigEndian.Uint64(counter)
}
