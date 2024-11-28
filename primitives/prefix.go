package primitives

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/raccoondb/errors"
	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/store"
	"github.com/sourcenetwork/raccoondb/types"
)

var ErrPrefixStore = errors.New("prefix store")

func newPrefixErr(method string, msg string, err error) error {
	return fmt.Errorf("%w: %v: %v: %w", ErrPrefixStore, method, msg, err)
}

func newKeyNilErr(method string) error {
	return fmt.Errorf("%w: %v: %w", ErrPrefixStore, method, store.ErrKeyNil)
}

func NewPrefixedKV(store store.KVStore, prefix []byte) store.KVStore {
	return &PrefixStore{
		store:  store,
		prefix: prefix,
	}
}

var _ iterator.BytesIterator = (*prefixStoreIterator)(nil)
var _ store.KVStore = (*PrefixStore)(nil)

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

func (kv *PrefixStore) Iterate(ctx context.Context, opt iterator.IteratorOpt) (iterator.BytesIterator, error) {
	iter, err := kv.store.Iterate(ctx, opt)
	if err != nil {
		return nil, newPrefixErr("Iterator", "creating iterator", err)
	}
	return &prefixStoreIterator{
		iter:   iterator.NewPrefixIterator(kv.prefix, iter),
		prefix: kv.prefix,
	}, nil
}

// prefixStoreIterator
type prefixStoreIterator struct {
	iter   *iterator.PrefixIterator[[]byte]
	prefix []byte
}

func (i *prefixStoreIterator) Finished() bool {
	return i.iter.Finished()
}

func (i *prefixStoreIterator) Next() error {
	return i.iter.Next()
}

// Key strips prefix from Key
func (i *prefixStoreIterator) CurrentKey() (key []byte) {
	key = i.iter.CurrentKey()
	if key == nil {
		return nil
	}
	return key[len(i.prefix):]
}

func (i *prefixStoreIterator) Value() types.Option[[]byte] {
	return i.iter.Value()
}

func (i *prefixStoreIterator) Close() error {
	return i.iter.Close()
}

func (i *prefixStoreIterator) GetParams() iterator.IteratorOpt {
	return i.iter.GetParams()
}
