package table

import (
	"context"
	"testing"

	"github.com/sourcenetwork/corekv/memory"
	"github.com/sourcenetwork/raccoondb/v2/iterator"
	"github.com/sourcenetwork/raccoondb/v2/marshal"
	"github.com/sourcenetwork/raccoondb/v2/store"
	"github.com/sourcenetwork/raccoondb/v2/store/corekv"
	"github.com/sourcenetwork/raccoondb/v2/store/test"
	"github.com/stretchr/testify/require"
)

type Record struct {
	Name string `json:"name"`
}

func Test_Table_Suite(t *testing.T) {
	factory := func() store.KVStore {
		m := marshal.BytesMarshaler{}
		kv := corekv.NewMemKV()
		t := NewTable(kv, m)
		return t
	}
	test.RunSuite(t, factory)
}

func TestIndexedObjectStore_Example(t *testing.T) {
	memcorekv := memory.NewDatastore(context.TODO())
	kv := corekv.WrapCoreKV(memcorekv)

	factory := func() Record { return Record{} }
	table := NewTable(kv, marshal.NewJSONMarshaler(factory))

	nameIdx, err := NewIndex(
		table,
		"name",
		func(record *Record) string { return record.Name },
		&marshal.StringMarshaler{},
	)
	require.NoError(t, err)

	ctx := context.TODO()
	created, err := table.Set(ctx, []byte("a"), Record{"bob"})
	require.NoError(t, err)
	require.True(t, bool(created))

	bucket := "bob"
	iter, err := nameIdx.IterateKeys(ctx, &bucket)
	require.NoError(t, err)
	test.DumpStore(t, kv)

	names, errs := iterator.Consume(ctx, iter)
	require.Empty(t, errs)

	require.Equal(t, []byte("a"), names[0])
}
