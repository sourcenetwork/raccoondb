// package raccoon is a lightweight and minimal graph store library
package raccoon

import (
    "fmt"
)


// RaccoonStore is a general purpose Graph store.
// It's built on top of a kv-store backend and provides common operations
// required while working with graphs.
//
// RaccoonStore has a lightweight secondary indexing framework,
// which can be used to optimize recurring access patterns.
//
// The main RaccoonStore type is a lightweight orchestrator on top
// of the primary index store and additional secondary index stores.
//
// TODO Make these atomic
type RaccoonStore[Edg Edge[N], N any] struct {
    indexes map[string]SecondaryIndexStore[Edg, N]
    recordStore EdgeStore[Edg, N]
}

// Build a Raccoon instance from a schema definition
func NewRaccoonStore[Edg Edge[N], N any](schema RaccoonSchema[Edg, N]) RaccoonStore[Edg, N]{
    idxs := make(map[string]SecondaryIndexStore[Edg, N])
    for _, idx := range schema.Indexes {
        store := NewSecondaryIdxStore(schema.Store, schema.KeysPrefix, idx.Name, schema.Keyer, idx.Mapper, schema.Marshaler)
        idxs[idx.Name] = store
    }
    primaryStore := NewEdgeStore(schema.Store, schema.KeysPrefix, schema.Keyer, schema.Marshaler)
    return RaccoonStore[Edg, N] {
        indexes: idxs,
        recordStore: primaryStore,
    }
}

// Set an Edge
func (s *RaccoonStore[Edg, N]) Set(edg Edg) error {
    for _, idx := range s.indexes {
        err := idx.Set(edg)
        if err != nil {
            return err
        }
    }
    return s.recordStore.Set(edg)
}

// Fetch and Edge from the source and destiniation nodes
func (s *RaccoonStore[Edg, N]) Get(source, dest N) (Option[Edg], error) {
    return s.recordStore.GetEdg(source, dest)
}

// Fetch all Edges from index idxName indexed by value
func (s *RaccoonStore[Edg, N]) GetByIdx(idxName string, value []byte) ([]Edg, error) {
    idx, ok := s.indexes[idxName]
    if !ok {
        return nil, fmt.Errorf("unknown index: %s", idxName)
    }

    return idx.GetByIdx(value)
}

// Filter and return all Edges from index idxName, indexed by value 
// that match the given predicate
func (s *RaccoonStore[Edg, N]) FilterByIdx(idxName string, value []byte, predicate Predicate[Edg]) ([]Edg, error) {
    idx, ok := s.indexes[idxName]
    if !ok {
        return nil, fmt.Errorf("unknown index: %s", idxName)
    }

    return idx.FilterByIdx(value, predicate)
}

// Return the direct ancestors of a Node,
// that is, returns all edges that have node as dest.
// Note: Does not recursively fetch Ancestors
func (s *RaccoonStore[Edg, N]) GetAncestors(node N) ([]Edg, error) { 
    return s.recordStore.GetAncestors(node)
}

// Return all direct sucessors of a Node,
// that is, returns all edges that have node as source.
// Note: Does not recursively fetch Ancestors
func (s *RaccoonStore[Edg, N]) GetSucessors(node N) ([]Edg, error) { 
    return s.recordStore.GetSucessors(node)
}

// Delete edg from the persistent storage
func (s *RaccoonStore[Edg, N]) Delete(edg Edg) (error) {
    for _, idx := range s.indexes {
        err := idx.Delete(edg)
        if err != nil {
            return err
        }
    }
    return s.recordStore.Delete(edg)
}

// TODO Define other methods such as filtering
