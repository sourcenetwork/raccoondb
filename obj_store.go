// package object_kvstore define a generic wrapper
// store for a protobuff serializable entity
package raccoon

type ObjectPredicate[Obj any] func(Obj) bool

type ObjectStore[Obj any] struct {
    store KVStore
    marshaler Marshaler[Obj]
    ider Ider[Obj]
}

func NewObjStore[O any](kv KVStore, marshaler Marshaler[O], ider Ider[O]) ObjectStore[O] {
    return ObjectStore[O]{
        store: kv,
        marshaler: marshaler,
        ider: ider,
    }
}

func (s *ObjectStore[Obj]) GetObject(id []byte) (Option[Obj], error) {
    bytes, err := s.store.Get(id)

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
    bytes, err := s.marshaler.Marshal(&obj)
    if err != nil {
        return err
    }

    key := s.ider.Id(obj)
    return s.store.Set(key, bytes)
}

func (s *ObjectStore[Obj]) DeleteById(id []byte) error {
    return s.store.Delete(id)
}

func (s *ObjectStore[Obj]) Delete(obj Obj) error {
    id := s.ider.Id(obj)
    return s.DeleteById(id)
}

func (s *ObjectStore[Obj]) HasById(id []byte) (bool, error) {
    return s.store.Has(id)
}

func (s *ObjectStore[Obj]) Has(obj Obj) (bool, error) {
    id := s.ider.Id(obj)
    return s.store.Has(id)
}

func (s *ObjectStore[Obj]) ListIds() ([][]byte, error) {
    iter := s.store.Iterator(nil, nil)

    var ids [][]byte

    for ; iter.Valid(); iter.Next() {
        key := iter.Key()
        if err := iter.Error(); err != nil {
            return nil, err
        }

        ids = append(ids, key)
    }

    return ids, nil
}

func (s *ObjectStore[Obj]) List() ([]Obj, error) {
    identity := func(o Obj) bool { return true }
    return s.Filter(identity)
}

func (s *ObjectStore[Obj]) Filter(predicate ObjectPredicate[Obj]) ([]Obj, error) {
    iter := s.store.Iterator(nil, nil)

    var objs []Obj

    for ; iter.Valid(); iter.Next() {
        if err := iter.Error(); err != nil {
            return nil, err
        }

        value := iter.Value()
        obj, err := s.marshaler.Unmarshal(value)
        if err != nil {
            return nil, err
        }
        if predicate(obj) {
            objs = append(objs, obj)
        }
    }

    return objs, nil
}
