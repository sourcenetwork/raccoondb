package primitives

import (
	"context"
	"testing"

	"github.com/sourcenetwork/raccoondb/store"
	"github.com/sourcenetwork/raccoondb/store/corekv"
	"github.com/sourcenetwork/raccoondb/store/test"
	"github.com/stretchr/testify/require"
)

var testKey = []byte("key")
var testVal = []byte("value")

func Test_CountedKV_StartsAtZero(t *testing.T) {
	kv := corekv.NewMemKV()
	ctx := context.TODO()
	ckv := NewCountedKVStore(kv)

	count, err := ckv.GetCount(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(0), count)
}

func Test_CountedKV_IncrementsWhenRecordIsSet(t *testing.T) {
	kv := corekv.NewMemKV()
	ctx := context.TODO()
	ckv := NewCountedKVStore(kv)

	_, err := ckv.Set(ctx, testKey, testVal)
	require.NoError(t, err)

	count, err := ckv.GetCount(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), count)
}

func Test_CountedKV_DecrementsAfterRecordIsRemoved(t *testing.T) {
	kv := corekv.NewMemKV()
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

func Test_CoutnedKV_Suite(t *testing.T) {
	factory := func() store.KVStore {
		kv := corekv.NewMemKV()
		ckv := NewCountedKVStore(kv)
		return ckv
	}
	test.RunSuite(t, factory)
}
