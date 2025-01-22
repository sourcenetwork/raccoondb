package table

import (
	"context"
	"testing"

	"github.com/sourcenetwork/raccoondb/v2/iterator"
	"github.com/sourcenetwork/raccoondb/v2/marshal"
	"github.com/sourcenetwork/raccoondb/v2/store"
	"github.com/sourcenetwork/raccoondb/v2/store/test"
	testutil "github.com/sourcenetwork/raccoondb/v2/test"
	"github.com/stretchr/testify/require"
)

var testKey []byte = []byte("test")

func setup(t *testing.T) (context.Context, *Table[Record], IndexReader[Record, string]) {
	kv := testutil.NewTestKV()

	factory := func() Record { return Record{} }
	table := NewTable(kv, marshal.NewJSONMarshaler(factory))

	nameIdx, err := NewIndex(
		table,
		"name",
		func(record *Record) string { return record.Name },
		&marshal.StringMarshaler{},
	)
	require.NoError(t, err)
	return context.TODO(), table, nameIdx
}

type Record struct {
	Name string `json:"name"`
}

func Test_Table_Suite(t *testing.T) {
	factory := func() store.KVStore {
		m := marshal.BytesMarshaler{}
		kv := testutil.NewTestKV()
		t := NewTable(kv, m)
		return t
	}
	test.RunSuite(t, factory)
}

func Test_Table_SettingObject_AddsItToIndex(t *testing.T) {
	ctx, table, idx := setup(t)

	created, err := table.Set(ctx, []byte("a"), Record{"bob"})
	require.NoError(t, err)
	require.True(t, bool(created))

	bucket := "bob"
	has, err := idx.Has(ctx, &bucket, []byte("a"))
	require.NoError(t, err)
	require.True(t, has)
}

func Test_Table_RemovingObjectRemovesItFromIdx(t *testing.T) {
	ctx, table, idx := setup(t)

	// given records 1 and 2 with name bob
	obj := Record{Name: "bob"}
	_, err := table.Set(ctx, []byte("1"), obj)
	require.NoError(t, err)
	_, err = table.Set(ctx, []byte("2"), obj)
	require.NoError(t, err)

	// when I delete record 1
	removed, err := table.Delete(ctx, []byte("2"))
	require.True(t, bool(removed))
	require.NoError(t, err)

	// then name index contains record 2 only
	bob := "bob"
	iter, err := idx.IterateKeys(ctx, &bob, store.NewOpenIterator())
	require.NoError(t, err)
	ids, errs := iterator.Consume(ctx, iter)
	require.Empty(t, errs)
	require.Equal(t, [][]byte{[]byte("1")}, ids)
	// and index does not have record 2
	has, err := idx.Has(ctx, &bob, []byte("2"))
	require.NoError(t, err)
	require.False(t, has)
}

func Test_Table_IterContainsIndexedObjects(t *testing.T) {
	ctx, table, idx := setup(t)

	table.Set(ctx, []byte("c"), Record{"alice"})
	table.Set(ctx, []byte("a"), Record{"bob"})
	table.Set(ctx, []byte("b"), Record{"bob"})

	bucket := "bob"
	iter, err := idx.IterateKeys(ctx, &bucket, store.NewOpenIterator())
	require.NoError(t, err)
	test.DumpStore(t, table.baseStore)

	keys, errs := iterator.Consume(ctx, iter)
	require.Empty(t, errs)
	want := [][]byte{
		[]byte("a"),
		[]byte("b"),
	}
	require.Equal(t, want, keys)
}

func Test_Table_IndexReturnsRightName(t *testing.T) {
	_, _, idx := setup(t)

	require.Equal(t, "name", idx.GetIndexName())
}

func Test_Table_IterateBuckets_ReturnsAllBuckets(t *testing.T) {
	ctx, table, idx := setup(t)
	table.Set(ctx, []byte("b"), Record{"bob"})
	table.Set(ctx, []byte("a"), Record{"alice"})

	iter, err := idx.IterateBuckets(ctx, NewOpenIterator[string]())
	require.NoError(t, err)

	buckets, errs := iterator.Consume(ctx, iter)
	require.Empty(t, errs)
	want := []string{
		"alice",
		"bob",
	}
	require.Equal(t, want, buckets)

	count, err := idx.GetBucketCount(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(2), count)
}

func Test_Table_UpdatingRecord_RecordMovesBuckets(t *testing.T) {
	ctx, table, idx := setup(t)
	table.Set(ctx, []byte("1"), Record{"alice"})
	table.Set(ctx, []byte("2"), Record{"alice"})

	// when I rename record 1 to bob
	created, err := table.Set(ctx, []byte("1"), Record{"bob"})
	require.False(t, bool(created))
	require.NoError(t, err)

	// Then bucket alice no longer contains record 1
	bkt := "alice"
	has, err := idx.Has(ctx, &bkt, []byte("1"))
	require.False(t, has)
	require.NoError(t, err)
	// and bucket bob contains record 1
	bkt = "bob"
	has, err = idx.Has(ctx, &bkt, []byte("1"))
	require.True(t, has)
	require.NoError(t, err)
}

func Test_Table_AddingIndexThenUpdating_BuildsIndexes(t *testing.T) {
	t.Skip()
	// Given table with record 1 and 2
	kv := testutil.NewTestKV()
	factory := func() Record { return Record{} }
	table := NewTable(kv, marshal.NewJSONMarshaler(factory))
	ctx := context.TODO()
	table.Set(ctx, []byte("1"), Record{"alice"})
	table.Set(ctx, []byte("2"), Record{"alice"})

	// When I add an Index and generated the indexes
	idx, err := NewIndex(
		table,
		"name",
		func(record *Record) string { return record.Name },
		&marshal.StringMarshaler{},
	)
	require.NoError(t, err)
	err = table.UpateIndexes(ctx)
	require.NoError(t, err)

	// Then idx has a buket for alice with keys 1 and 2
	bucket := "alice"
	has, err := idx.Has(ctx, &bucket, []byte("1"))
	require.NoError(t, err)
	require.True(t, has)
	idx.Has(ctx, &bucket, []byte("2"))
	require.NoError(t, err)
	require.True(t, has)
}
