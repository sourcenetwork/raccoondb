package raccoon

const (
    inPrefix = "/in"
    outPrefix = "/out"
)

// TODO refactor type constraints to DirectedEdge

// Store Wrapper for records indexed by a certain field
// Keys in store have the pattern: {prefix}/{indexName}/{mappedValue}/{edgeKey}
type SecondaryIndexStore[Edg Edge[N], N any] struct {
    store KVStore
    prefix []byte
    indexName string
    mapper Mapper[Edg]
    keyer SecondaryIdxKeyer[Edg, N]
    marshaler Marshaler[Edg]
}

func NewSecondaryIdxStore[Edg Edge[N], N any](store KVStore, prefix []byte, indexName string, keyer NodeKeyer[N], mapper Mapper[Edg], marshaler Marshaler[Edg]) SecondaryIndexStore[Edg, N] {
    idxKeyer := NewSecIdxKeyer[Edg, N](keyer, mapper, indexName, prefix)
    return SecondaryIndexStore[Edg, N] {
        store: store,
        prefix: prefix,
        indexName: indexName,
        mapper: mapper,
        keyer: idxKeyer,
        marshaler: marshaler,
    }
}

func (s *SecondaryIndexStore[Edg, N]) Set(edg Edg) error {
    bytes, err := s.marshaler.Marshal(&edg)
    if err != nil {
        return err
    }

    key := s.keyer.Key(edg)
    s.store.Set(key.ToBytes(), bytes)
    return nil
}


func (s *SecondaryIndexStore[Edg, N]) Delete(edg Edg) error {
    key := s.keyer.Key(edg)
    s.store.Delete(key.ToBytes())
    return nil
}

func (s *SecondaryIndexStore[Edg, N]) GetByIdx(idx []byte) ([]Edg, error) {
    min, max := s.keyer.IterKeys(idx)
    iter := s.store.Iterator(min.ToBytes(), max.ToBytes())
    return consumeIter[Edg](iter, s.marshaler)
}

func (s *SecondaryIndexStore[Edg, N]) FilterByIdx(idx []byte, predicate Predicate[Edg]) ([]Edg, error) {
    min, max := s.keyer.IterKeys(idx)
    iter := s.store.Iterator(min.ToBytes(), max.ToBytes())
    return filterConsumeIter(iter, s.marshaler, predicate)
}


// Store for Graph Edges
// 
// Provides the usual Get, Set, Delete operations for records.
// Also exposes a higher level operation to return all records which
// contain a Node's sucessor and ancestor.
//
// EdgeStore maintains an internal index for a node's ancestors and sucessors,
// allowing for fast lookup.
// Setting and Deleting records updates internal indexes to guarantee consistency.
// TODO make it atomic
type EdgeStore[Edg Edge[N], N any] struct {
    store KVStore
    prefix []byte
    keyer EdgeKeyer[N]
    marshaler Marshaler[Edg]
}

func NewEdgeStore[Edg Edge[N], N any](store KVStore, prefix []byte, keyer NodeKeyer[N], marshaler Marshaler[Edg]) EdgeStore[Edg, N] {
    base := Key{}.Append(prefix)
    edgKeyer := NewEdgeKeyer[N](keyer, base, []byte(inPrefix), []byte(outPrefix))
    return EdgeStore[Edg, N] {
        store: store,
        prefix: prefix,
        keyer: edgKeyer,
        marshaler: marshaler,
    }
}

// Set a record in the backend storage engine.
func (s *EdgeStore[Edg, N]) Set(edg Edg) error {
    bytes, err := s.marshaler.Marshal(&edg)
    if err != nil {
        return err
    }

    src := edg.GetSource()
    dst := edg.GetDest()

    outKey := s.keyer.OutgoingKey(src, dst)
    s.store.Set(outKey.ToBytes(), bytes)

    inKey := s.keyer.IncomingKey(src, dst)
    s.store.Set(inKey.ToBytes(), bytes)

    return nil
}

// Fetch stored Record for the given edge.
// Return record, an ok flag and a potential error
// If the OK flag is false, no record was found or an error ocurred
func (s *EdgeStore[Edg, N]) GetEdg(source, dest N) (Option[Edg], error) {
    key := s.keyer.OutgoingKey(source, dest)

    bytes, err := s.store.Get(key.ToBytes())
    if bytes == nil || err != nil {
        return None[Edg](), nil
    }

    edg, err := s.marshaler.Unmarshal(bytes)
    if err != nil {
        return None[Edg](), err
    }
    
    return Some(edg), nil
}

// Delete record associated to edge entry from store
func (s *EdgeStore[Edg, N]) Delete(edg Edg) error {
    source := edg.GetSource()
    dest := edg.GetDest()

    key := s.keyer.IncomingKey(source, dest)
    s.store.Delete(key.ToBytes())

    key = s.keyer.OutgoingKey(source, dest)
    s.store.Delete(key.ToBytes())

    return nil
}

// Return all Records in which `node` is the edge's source
func (s *EdgeStore[Edg, N]) GetSucessors(node N) ([]Edg, error) {
    min, max := s.keyer.SucessorsIterKeys(node)
    iter := s.store.Iterator(min.ToBytes(), max.ToBytes())
    return consumeIter[Edg](iter, s.marshaler)
}

// Return all Records in which `node` is the edge's target
func (s *EdgeStore[Edg, N]) GetAncestors(node N) ([]Edg, error) {
    min, max := s.keyer.AncestorsIterKey(node)
    iter := s.store.Iterator(min.ToBytes(), max.ToBytes())
    return consumeIter[Edg](iter, s.marshaler)
}

// TODO define FilterAncestors and FilterSucessors
