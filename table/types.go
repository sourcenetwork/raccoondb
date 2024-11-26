package table

import "github.com/sourcenetwork/raccoondb/iterator"

// ObjKeyIter represents an iterator whose values are object keys
type ObjKeyIter iterator.BytesIterator

type Catalogue struct {
	ObjectCount uint64
	IndexesData map[string]IndexData
}

type IndexData struct {
	Name               string
	BucketCount        uint64
	IndexedObjectCount uint64
}

type AutoIncermenterMutations interface {
}

type TableMutations interface {
}
