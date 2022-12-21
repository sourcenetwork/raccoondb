// package object_kvstore define a generic wrapper
// store for a protobuff serializable entity
package raccoon


type ObjectStore[Obj any] struct {
    store KVStore
    prefix []byte
    marshaler Marshaler[Obj]
    ider Ider[Obj]
    baseKey Key
}

func NewObjStore[O any](kv KVStore, prefix []byte, marshaler Marshaler[O], ider Ider[O]) ObjectStore[O] {
    return ObjectStore[O]{
        store: kv,
        prefix: prefix,
        marshaler: marshaler,
        ider: ider,
        baseKey: Key{}.Append(prefix),
    }
}

func (s *ObjectStore[Obj]) GetObject(id []byte) (Option[Obj], error) {
    key := s.baseKey.Append(id)
    bytes, err := s.store.Get(key.ToBytes())

    if err != nil || bytes == nil{
        return None[Obj](), err
    }

    obj, err := s.marshaler.Unmarshal(bytes)
    if err != nil {
        return None[Obj](), err
    }


    return Some(obj), nil
}

func (s *ObjectStore[Obj]) SetObject(obj Obj) error {
    bytes, err := s.marshaler.Marshal(obj)
    if err != nil {
        return err
    }

    id := s.ider.Id(obj)
    key := s.baseKey.Append(id)
    s.store.Set(key.ToBytes(), bytes)
    return nil
}

func (s *ObjectStore[Obj]) DeleteById(id []byte) error {
    key := s.baseKey.Append(id)
    s.store.Delete(key.ToBytes())
    return nil
}

func (s *ObjectStore[Obj]) Delete(obj Obj) error {
    id := s.ider.Id(obj)
    return s.DeleteById(id)
}
