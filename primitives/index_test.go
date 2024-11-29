package primitives

import (
	"context"
	"testing"

	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/store/corekv"
	"github.com/sourcenetwork/raccoondb/store/test"
	"github.com/stretchr/testify/require"
)

var testBucket []byte = []byte("bucket")

func Test_FieldIndexStore_ValesInBucket_CanIter(t *testing.T) {
	ctx := context.TODO()
	kv := corekv.NewMemKV()
	idx := NewFieldIndexStore(kv)

	created, err := idx.IndexValue(ctx, testBucket, []byte("v1"))
	require.True(t, bool(created))
	require.NoError(t, err)

	created, err = idx.IndexValue(ctx, testBucket, []byte("v2"))
	require.True(t, bool(created))
	require.NoError(t, err)

	iter, err := idx.IterateBucketItems(ctx, testBucket)
	require.NoError(t, err)

	pairs := iterator.ConsumePairs(ctx, iter)
	want := []iterator.Pair[[]byte]{
		iterator.NewPair([]byte("v1"), []byte("v1")),
		iterator.NewPair([]byte("v2"), []byte("v2")),
	}
	require.Equal(t, want, pairs)
}

func Test_FieldIndexStore_GetBucketCount_ReturnsCount(t *testing.T) {
	ctx := context.TODO()
	kv := corekv.NewMemKV()
	idx := NewFieldIndexStore(kv)

	created, err := idx.IndexValue(ctx, []byte("b1"), []byte("v1"))
	require.True(t, bool(created))
	require.NoError(t, err)

	created, err = idx.IndexValue(ctx, []byte("b2"), []byte("v2"))
	require.True(t, bool(created))
	require.NoError(t, err)

	count, err := idx.GetBucketCount(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(2), count)
}

func Test_FieldIndexStore_RemovingLastItemFromBucket_DeletesAndDecrementBucket(t *testing.T) {
	ctx := context.TODO()
	kv := corekv.NewMemKV()
	idx := NewFieldIndexStore(kv)

	created, err := idx.IndexValue(ctx, []byte("b1"), []byte("v1"))
	require.True(t, bool(created))
	require.NoError(t, err)

	created, err = idx.IndexValue(ctx, []byte("b2"), []byte("v2"))
	require.True(t, bool(created))
	require.NoError(t, err)

	removed, err := idx.RemoveItem(ctx, []byte("b2"), []byte("v2"))
	require.True(t, bool(removed))
	require.NoError(t, err)

	test.DumpStore(t, kv)

	count, err := idx.GetBucketCount(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), count)

	bucketsIter, err := idx.IterateBuckets(ctx)
	require.NoError(t, err)
	buckets, errs := iterator.Consume(ctx, bucketsIter)
	require.Empty(t, errs)
	want := [][]byte{
		[]byte("b1"),
	}
	require.Equal(t, want, buckets)
}
