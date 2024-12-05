package primitives

import (
	"context"
	"testing"

	"github.com/sourcenetwork/raccoondb/v2/store"
	"github.com/sourcenetwork/raccoondb/v2/store/corekv"
	"github.com/sourcenetwork/raccoondb/v2/store/test"
	"github.com/stretchr/testify/require"
)

func Test_PrefixKV_Suite(t *testing.T) {
	prefix := []byte("prefix/")
	factory := func() store.KVStore {
		kv := corekv.NewMemKV()
		return NewPrefixedKV(kv, prefix)
	}
	test.RunSuite(t, factory)
}

func Test_PrefixKV_AddsPrefixToAllElems(t *testing.T) {
	prefix := []byte("prefix/")
	kv := corekv.NewMemKV()
	pkv := NewPrefixedKV(kv, prefix)

	ctx := context.TODO()
	_, err := pkv.Set(ctx, testKey, testVal)
	require.NoError(t, err)

	opt, err := kv.Get(ctx, concatKey(prefix, testKey))
	require.NoError(t, err)
	require.False(t, opt.Empty())
	require.Equal(t, testVal, opt.GetValue())
}

func Test_PrefixKV_BaseStoreWithElementsAfterPrefixAreNotIncludedInIterator(t *testing.T) {
	// Given parent store with xtest entry
	ctx := context.TODO()
	kv := corekv.NewMemKV()
	key := []byte("xtest")
	_, err := kv.Set(ctx, key, key)
	require.NoError(t, err)

	// given prefixe store with testKey entry
	prefix := []byte("prefix/")
	pkv := NewPrefixedKV(kv, prefix)
	_, err = pkv.Set(ctx, testKey, testVal)
	require.NoError(t, err)

	// when I iterate all items in prefix store
	iter, err := pkv.Iterate(ctx, store.NewOpenIterator())
	require.NoError(t, err)

	// then only the entry for testKey is returned
	err = iter.Next(ctx)
	require.NoError(t, err)
	require.Equal(t, testKey, iter.CurrentKey())

	err = iter.Next(ctx)
	require.NoError(t, err)
	require.True(t, iter.Finished())
}
