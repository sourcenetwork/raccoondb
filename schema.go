package raccoon

// Secondary Index schema definition
// Specifies the index name and a mapper function.
type SecondaryIndex[Edg Edge[N], N any] struct {
    Name string
    Mapper Mapper[Edg]
}


// Raccoon Store schema
type RaccoonSchema[Edg Edge[N], N any] struct {
    Indexes []SecondaryIndex[Edg, N]
    Store KVStore
    KeysPrefix []byte
    Keyer NodeKeyer[N]
    Marshaler Marshaler[Edg]
}
