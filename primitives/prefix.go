package primitives

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/raccoondb/v2/errors"
	"github.com/sourcenetwork/raccoondb/v2/store"
	"github.com/sourcenetwork/raccoondb/v2/types"
)

var _ store.KVStore = (*PrefixStore)(nil)

// ErrPrefixStore is a top level error for all errors produced by PrefixStore
var ErrPrefixStore = errors.New("prefix store")

func newPrefixErr(method string, msg string, err error) error {
	return fmt.Errorf("%w: %v: %v: %w", ErrPrefixStore, method, msg, err)
}

func newKeyNilErr(method string) error {
	return fmt.Errorf("%w: %v: %w", ErrPrefixStore, method, store.ErrKeyNil)
}

// NewPrefixedKV returns a KVStore which wraps store to automatically add prefix
// to all keys.
// Prefixed KVStores can be used to scope or add namespaces to a KVStore
func NewPrefixedKV(store store.KVStore, prefix []byte) store.KVStore {
	return &PrefixStore{
		store:  store,
		prefix: prefix,
	}
}

// PrefixStore implements raccoon's store.KVStore to a store.KVStore by wrapping its methods with a global prefix
type PrefixStore struct {
	store  store.KVStore
	prefix []byte
}

func (kv *PrefixStore) joinKey(key []byte) []byte { return concatKey(kv.prefix, key) }

func (kv *PrefixStore) Get(ctx context.Context, key []byte) (types.Option[[]byte], error) {
	if key == nil {
		return types.None[[]byte](), newKeyNilErr("Get")
	}

	key = kv.joinKey(key)
	opt, err := kv.store.Get(ctx, key)
	if err != nil {
		return types.None[[]byte](), newPrefixErr("Get", "fetching key", err)
	}
	return opt, nil
}

func (kv *PrefixStore) Has(ctx context.Context, key []byte) (bool, error) {
	if key == nil {
		return false, newKeyNilErr("Get")
	}

	key = kv.joinKey(key)
	has, err := kv.store.Has(ctx, key)
	if err != nil {
		return false, newPrefixErr("Has", "checking", err)
	}
	return has, nil
}

func (kv *PrefixStore) Set(ctx context.Context, key, value []byte) (store.KeyCreated, error) {
	if key == nil {
		return false, newKeyNilErr("Get")
	}

	key = kv.joinKey(key)
	created, err := kv.store.Set(ctx, key, value)
	if err != nil {
		return false, newPrefixErr("Set", "setting value", err)
	}
	return created, nil
}

func (kv *PrefixStore) Delete(ctx context.Context, key []byte) (store.KeyRemoved, error) {
	if key == nil {
		return false, newKeyNilErr("Get")
	}

	key = kv.joinKey(key)
	deleted, err := kv.store.Delete(ctx, key)
	if err != nil {
		return false, newPrefixErr("Delete", "removing record", err)
	}
	return deleted, nil
}

func (kv *PrefixStore) Iterate(ctx context.Context, opt store.IteratorOpt) (store.StoreIterator[[]byte], error) {
	iter, err := store.IteratePrefix(ctx, kv.store, kv.prefix, opt, true)
	if err != nil {
		return nil, newPrefixErr("Iterator", "creating iterator", err)
	}
	return store.ToStoreIter(iter, opt), nil
}
