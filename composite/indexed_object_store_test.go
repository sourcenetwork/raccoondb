package composite

import (
	"testing"

	"github.com/sourcenetwork/raccoondb/marshal"
)

type Record struct {
	Name string `json:name`
}

func TestIndexedObjectStore_Example(t *testing.T) {
	nameIdx := Index[Record, string]{
		Name: "name",
		Extractor: func(record *Record) string {
			return record.Name
		},
		Marshaler: marshal.StringMarshaler{},
	}
	factory := func() Record {
		return Record{}
	}

	schema := &IndexedObjectStoreSchema[Record]{
		Marshaler: marshal.NewJSONMarshaler(factory),
		Indexes: []*Index[Record, any]{
			&nameIdx,
		},
	}

}
