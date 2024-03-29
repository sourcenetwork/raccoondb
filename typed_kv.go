package raccoon


import (
    "fmt"
)


var _ ObjKV[any] = (*objKV[any])(nil)


// Return an ObjKV from a KVStore using marshaler to (un)marshal objects.
// Prefixes all keys with keyPrefix.
func NewObjKV[O any](kv KVStore, keyPrefix []byte, marshaler Marshaler[O]) ObjKV[O] {
    return &objKV[O]{
        store: kv,
        prefix: keyPrefix,
        marshaler: marshaler,
        baseKey: Key{}.Append(keyPrefix),
    }
}


// objKV implements raccoon's ObjKV interface
type objKV[Obj any] struct {
    store KVStore
    prefix []byte
    marshaler Marshaler[Obj]
    baseKey Key
}

// Fetch object from store using the given key
func (s *objKV[Obj]) Get(key []byte) (Option[Obj], error) {
    key = s.getFullKey(key)

    bytes, err := s.store.Get(key)
    if err != nil || bytes == nil{
        err = fmt.Errorf("failed to fetch key %v: %w", key, err)
        return None[Obj](), err
    }

    obj, err := s.marshaler.Unmarshal(bytes)
    if err != nil {
        err = fmt.Errorf("failed unmarshaling obj from key %v: %w", key, err)
        return None[Obj](), err
    }

    return Some(obj), nil
}

// Set key with obj
func (s *objKV[Obj]) Set(key []byte, obj Obj) error {
    key = s.getFullKey(key)

    bytes, err := s.marshaler.Marshal(&obj)
    if err != nil {
        return fmt.Errorf("failed unmarshaling obj %v: %w", key, err)
    }

    return s.store.Set(key, bytes)
}

// Remove key from store
func (s *objKV[Obj]) Delete(key []byte) error {
    key = s.getFullKey(key)

    err := s.store.Delete(key)
    if err != nil {
        return fmt.Errorf("failed deleting obj %v: %w", key, err)
    }
    return nil
}

// Check whether key exists in KVStore
func (s *objKV[Obj]) Has(key []byte) (bool, error) {
    key = s.getFullKey(key)

    has, err := s.store.Has(key)
    if err != nil {
        return false, fmt.Errorf("failed checking for obj %v: %w", key, err)
    }
    return has, nil
}

// Return bye slice generated by appending key to the store prefix
func (s *objKV[Obj]) getFullKey(key []byte) []byte {
    return s.baseKey.Append(key).ToBytes()
}
