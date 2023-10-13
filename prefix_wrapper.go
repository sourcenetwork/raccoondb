package raccoon

import (
    "bytes"
)

// WrapperKV implements raccoon's KVStore to a KVStore by wrapping its methods with a global prefix
type WrapperKV struct {
    store KVStore
    prefix Key
}

func NewWrapperKV(store KVStore, prefix []byte) KVStore {
    prefixKey := Key{}.Append(prefix)
    return &WrapperKV{
        store: store,
        prefix: prefixKey,
    }
}

func (kv *WrapperKV) Get(key []byte) ([]byte, error) {
    key = kv.prefix.Append(key).ToBytes()
    return kv.store.Get(key)
}

func (kv *WrapperKV) Has(key []byte) (bool, error) {
    key = kv.prefix.Append(key).ToBytes()
    return kv.store.Has(key)
}

func (kv *WrapperKV) Set(key, value []byte) error {
    key = kv.prefix.Append(key).ToBytes()
    return kv.store.Set(key, value)
}

func (kv *WrapperKV) Delete(key []byte) error {
    key = kv.prefix.Append(key).ToBytes()
    return kv.store.Delete(key)
}


func (kv *WrapperKV) Iterator(start, end []byte) Iterator {
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
        iter: iter,
        done: false,
    }
}

type prefixIterator struct {
    prefix []byte
    iter Iterator
    done bool
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
