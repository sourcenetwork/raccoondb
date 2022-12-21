package raccoon


// Type alias for KVStore
type KVStore interface {
    Get(key []byte) ([]byte, error)
    Has(key []byte) (bool, error)
    Set(key, value []byte) error
    Delete(key []byte) error
    Iterator(start, end []byte) Iterator
}

// Type ObjKV wraps KVStore by abstracting object marshaling
type ObjKV[T any] interface {
    // Fetch object from store using the given key
    Get(key []byte) (Option[T], error)

    // Check whether key exists in KVStore
    Has(key []byte) (bool, error)

    // Set key with obj
    Set(key []byte, obj T) error

    // Remove key from store
    Delete(key []byte) error
}

type Iterator interface {
    // Valid returns whether the current iterator is valid. Once invalid, the Iterator remains
    // invalid forever.
    Valid() bool

    // Next moves the iterator to the next key in the database, as defined by order of iteration.
    // If Valid returns false, this method will panic.
    Next()

    // Key returns the key at the current position. Panics if the iterator is invalid.
    // CONTRACT: key readonly []byte
    Key() (key []byte)

    // Value returns the value at the current position. Panics if the iterator is invalid.
    // CONTRACT: value readonly []byte
    Value() (value []byte)

    // Error returns the last error encountered by the iterator, if any.
    Error() error

    // Close closes the iterator, relasing any allocated resources.
    Close() error
}


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

type TypePointer[T any] interface {
    *T
}
