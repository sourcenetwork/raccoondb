package primitives

import (
	"context"
	"testing"

	"github.com/sourcenetwork/raccoondb/v2/store"
	"github.com/sourcenetwork/raccoondb/v2/store/test"
	testutil "github.com/sourcenetwork/raccoondb/v2/test"
	"github.com/stretchr/testify/require"
)

var testKey = []byte("key")
var testVal = []byte("value")

func Test_CountedKV_StartsAtZero(t *testing.T) {
	kv := testutil.NewTestKV()
	ctx := context.TODO()
	ckv := NewCountedKVStore(kv)

	count, err := ckv.GetCount(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(0), count)
}

func Test_CountedKV_IncrementsWhenRecordIsSet(t *testing.T) {
	kv := testutil.NewTestKV()
	ctx := context.TODO()
	ckv := NewCountedKVStore(kv)

	_, err := ckv.Set(ctx, testKey, testVal)
	require.NoError(t, err)

	count, err := ckv.GetCount(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), count)
}

func Test_CountedKV_DecrementsAfterRecordIsRemoved(t *testing.T) {
	kv := testutil.NewTestKV()
	ctx := context.TODO()
	ckv := NewCountedKVStore(kv)

	_, err := ckv.Set(ctx, testKey, testVal)
	require.NoError(t, err)

	_, err = ckv.Delete(ctx, testKey)
	require.NoError(t, err)

	count, err := ckv.GetCount(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(0), count)
}

func Test_CountedKV_Suite(t *testing.T) {
	factory := func() store.KVStore {
		kv := testutil.NewTestKV()
		ckv := NewCountedKVStore(kv)
		return ckv
	}
	test.RunSuite(t, factory)
}
