package raccoon

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
    start = kv.prefix.Append(start).ToBytes()
    end = kv.prefix.Append(end).ToBytes()
    return kv.store.Iterator(start, end)
}
