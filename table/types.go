package table

import "github.com/sourcenetwork/raccoondb/v2/iterator"

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
