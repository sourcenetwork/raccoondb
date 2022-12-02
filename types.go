package raccoon

import (
    "github.com/cosmos/cosmos-sdk/store/types"
)

// Type alias for KVStore
type KVStore types.KVStore

// Predicate type defines a filter over a Record set
type Predicate[T any] func(T) bool

// Mapper type is used to map a Record to a value used by secundary indexes
type Mapper[T any] func(T) []byte


// Marshaler marshalls and unmarshalls objects into a byte array,
// which gets persisted in the underlying store
type Marshaler[T any] interface {
    Marshal(*T) ([]byte, error)
    Unmarshal([]byte) (T, error)
}

// NodeKeyer interface specify methods to map a node into an
// identifier key.
// Keys are an user defined byte array.
// Key values are opaque to the library however they must be deterministic and unique.
// Furthermore, keys should have known bounds, which are used for
// iteration.
// 
// NodeKeyer decouples node identification from data, which allows
// for different key generation strategies based on usage patterns.
type NodeKeyer[T any] interface {
    // Maps a Node into a key.
    Key(T) []byte

    // Return the lowest possible key
    MinKey() []byte

    // Return the highest possible key
    MaxKey() []byte
}

// Return a unique key (id) for the given object
type Ider[T any] interface {
    Id(T) []byte
}


// TODO Fix edge types

// Edge represents a protobuff serializable type which is persisted.
// Edge contains a source and a target Node.
// Users implementing the Edge interface are free to add extra 
// data to the edge type.
// Edges are uniquely identified through their nodes
type Edge[Node any] interface {
    GetSource() Node
    GetDest() Node
    // NOTE Zanzibar's data model also stores the full node in each edge.
    // This has pros and cons.
    // Cons is duplicating data and more storage requirements
}

// DirectedEdge specialization of an Edge.
type DirectedEdge[Node any] interface {
    Edge[Node]
}

// UndirectedEdge specialization of an Edge.
type UndirectedEdge[Node any] interface {
    Edge[Node]
}
