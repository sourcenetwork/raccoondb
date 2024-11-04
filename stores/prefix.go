package stores

func NewPrefixedKV(store KVStore, prefix []byte) KVStore {
	return nil
}

/*
import (
	"bytes"
	"context"

	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/types"
)

var _ iterator.BytesIterator = (*prefixIterator)(nil)
var _ KVStore = (*PrefixStore)(nil)

// PrefixStore implements raccoon's KVStore to a KVStore by wrapping its methods with a global prefix
type PrefixStore struct {
	store  KVStore
	prefix []byte
}

func NewPrefixedKV(store KVStore, prefix []byte) KVStore {
	prefixKey := Key{}.Append(prefix)
	return &PrefixStore{
		store:  store,
		prefix: prefixKey,
	}
}

func (kv *PrefixStore) Get(ctx context.Context, key []byte) ([]byte, error) {
	key = kv.prefix.Append(key).ToBytes()
	return kv.store.Get(ctx, key)
}

func (kv *PrefixStore) Has(ctx context.Context, key []byte) (bool, error) {
	key = kv.prefix.Append(key).ToBytes()
	return kv.store.Has(key)
}

func (kv *PrefixStore) Set(ctx context.Context, key, value []byte) error {
	key = kv.prefix.Append(key).ToBytes()
	return kv.store.Set(key, value)
}

func (kv *PrefixStore) Delete(ctx context.Context, key []byte) error {
	key = kv.prefix.Append(key).ToBytes()
	return kv.store.Delete(key)
}

func (kv *PrefixStore) Iterator(ctx context.Context, start, end []byte) iterator.BytesIterator {
	return newPrefixIterator(kv.prefix, start, end, kv.store)
}

func newPrefixIterator(prefix Key, start, end []byte, store KVStore) *prefixIterator {
	prefixBytes := prefix.Append(nil).ToBytes()
	start = prefix.Append(start).ToBytes()
	// if end is nil, the iterator must be unbounded
	if end != nil {
		end = prefix.Append(end).ToBytes()
	}

	iter := store.Iterator(start, end)

	return &prefixIterator{
		prefix: prefixBytes,
		iter:   iter,
		done:   false,
	}
}

type prefixIterator struct {
	prefix []byte
	iter   types.BytesIterator
	done   bool
}

func (i *prefixIterator) Valid() bool {
	if i.done {
		return false
	}
	return i.iter.Valid()
}

// Next steps the iterator to the next value
// if the next value does not contain prefix, the scan is done
func (i *prefixIterator) Next() {
	i.iter.Next()

	if !i.iter.Valid() || !bytes.HasPrefix(i.iter.Key(), i.prefix) {
		i.done = true
	}
}

// Key strips prefix from Key
func (i *prefixIterator) Key() (key []byte) {
	key = i.iter.Key()
	return key[len(i.prefix):]
}

func (i *prefixIterator) Value() (value []byte) {
	return i.iter.Value()
}

func (i *prefixIterator) Error() error {
	return i.iter.Error()
}

func (i *prefixIterator) Close() error {
	return i.iter.Close()
}

*/
