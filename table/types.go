package table

import "github.com/sourcenetwork/raccoondb/v2/iterator"

// ObjKeyIter represents an iterator whose values are object keys
type ObjKeyIter iterator.BytesIterator

// Catalogue models metadata tracked by a table
type Catalogue struct {
	ObjectCount uint64
	IndexesData map[string]IndexData
}

// IndexData models metadata tracked by an Index
type IndexData struct {
	Name               string
	BucketCount        uint64
	IndexedObjectCount uint64
}
