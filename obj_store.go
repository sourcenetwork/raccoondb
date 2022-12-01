// package object_kvstore define a generic wrapper
// store for a protobuff serializable entity
package raccoon

import (
    cosmos "github.com/cosmos/cosmos-sdk/store/types"

)


type ObjectStore[Obj any] struct {
    store cosmos.KVStore
    prefix []byte
    marshaller Marshaler[Obj]
    ider Ider[Obj]
    baseKey Key
}

func (s *ObjectStore[Obj]) GetObject(id []byte) (Option[Obj], error) {
    key := s.baseKey.Append(id)
    bytes := s.store.Get(key.ToBytes())
    if bytes == nil {
        return None[Obj](), nil
    }

    obj, err := s.marshaller.Unmarshal(bytes)
    if err != nil {
        return None[Obj](), err
    }


    return Some(obj), nil
}

func (s *ObjectStore[Obj]) SetObject(obj Obj) error {
    bytes, err := s.marshaller.Marshal(obj)
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

