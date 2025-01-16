package test

import (
	"context"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/sourcenetwork/raccoondb/v2/store"
	"github.com/stretchr/testify/require"
)

var testKey = []byte("key")
var testVal = []byte("value")

func test_SettingThenGet_ReturnsDoc(t *testing.T, kv store.KVStore) {
	ctx := context.TODO()

	_, err := kv.Set(ctx, testKey, testVal)
	require.NoError(t, err)

	opt, err := kv.Get(ctx, testKey)
	require.NoError(t, err)
	require.False(t, opt.Empty())
	require.Equal(t, testVal, opt.GetValue())
}

func test_GetDocumentNotRegistered_ReturnsEmptyOption(t *testing.T, kv store.KVStore) {
	ctx := context.TODO()

	opt, err := kv.Get(ctx, testKey)
	require.NoError(t, err)
	require.True(t, opt.Empty())
}

func test_GetWithNilKey_Errors(t *testing.T, kv store.KVStore) {
	ctx := context.TODO()

	opt, err := kv.Get(ctx, nil)
	require.ErrorIs(t, err, store.ErrKeyNil)
	require.True(t, opt.Empty())
}

func test_SetNewValue_ReturnsCreatedTrue(t *testing.T, kv store.KVStore) {
	ctx := context.TODO()

	created, err := kv.Set(ctx, testKey, testVal)
	require.NoError(t, err)
	require.True(t, bool(created))
}

func test_UpdateKey_ReturnsCreatedFalse(t *testing.T, kv store.KVStore) {
	ctx := context.TODO()

	_, err := kv.Set(ctx, testKey, testVal)
	require.NoError(t, err)

	created, err := kv.Set(ctx, testKey, testVal)
	require.NoError(t, err)
	require.False(t, bool(created))
}

func test_DeleteKeyNotSet_ReturnsRemovedFalse(t *testing.T, kv store.KVStore) {
	ctx := context.TODO()

	removed, err := kv.Delete(ctx, testKey)
	require.NoError(t, err)
	require.False(t, bool(removed))
}

func test_SetThenDelete_KeyNoLongerInStore(t *testing.T, kv store.KVStore) {
	ctx := context.TODO()

	_, err := kv.Set(ctx, testKey, testVal)
	require.NoError(t, err)

	_, err = kv.Delete(ctx, testKey)
	require.NoError(t, err)

	opt, err := kv.Get(ctx, testKey)
	require.NoError(t, err)
	require.True(t, opt.Empty())
}

func test_SetNilKey_Errors(t *testing.T, kv store.KVStore) {
	ctx := context.TODO()

	created, err := kv.Set(ctx, nil, testVal)
	require.ErrorIs(t, err, store.ErrKeyNil)
	require.False(t, bool(created))
}

func test_SetNilValue_Ok(t *testing.T, kv store.KVStore) {
	t.Skip()
	// this is an annoying edge case
	// comet does not supoprt this.
	// wrapping it in order for it to support would require an option
	// at the storage level which would solve the problem but it's overkill
	// for an edge case
	ctx := context.TODO()

	created, err := kv.Set(ctx, testKey, nil)
	require.NoError(t, err)
	require.True(t, bool(created))
}

func test_SetNilKV_ReturnsNilKeyErr(t *testing.T, kv store.KVStore) {
	ctx := context.TODO()

	created, err := kv.Set(ctx, nil, nil)
	require.ErrorIs(t, err, store.ErrKeyNil)
	require.False(t, bool(created))
}

func test_DeleteNilKey_Errors(t *testing.T, kv store.KVStore) {
	ctx := context.TODO()

	created, err := kv.Delete(ctx, nil)
	require.ErrorIs(t, err, store.ErrKeyNil)
	require.False(t, bool(created))
}

func test_HasNilKey_Errors(t *testing.T, kv store.KVStore) {
	ctx := context.TODO()

	has, err := kv.Has(ctx, nil)
	require.ErrorIs(t, err, store.ErrKeyNil)
	require.False(t, bool(has))
}

func test_Has_TrueWhenSet(t *testing.T, kv store.KVStore) {
	ctx := context.TODO()

	_, err := kv.Set(ctx, testKey, testVal)
	require.NoError(t, err)

	has, err := kv.Has(ctx, testKey)
	require.True(t, has)
	require.NoError(t, err)
}

func test_Has_FalseWhenNotSet(t *testing.T, kv store.KVStore) {
	ctx := context.TODO()

	has, err := kv.Has(ctx, testKey)
	require.False(t, has)
	require.NoError(t, err)
}

func test_Iterate_ReturnsIteratorOverAllItems(t *testing.T, kv store.KVStore) {
	ctx := context.TODO()
	testData := []string{
		"abc",
		"dave",
		"potato",
	}
	for _, val := range testData {
		_, err := kv.Set(ctx, []byte(val), []byte(val))
		require.NoError(t, err)
	}

	iter, err := kv.Iterate(ctx, store.NewOpenIterator())
	require.NoError(t, err)

	require.Nil(t, iter.CurrentKey())
	opt := iter.Value()
	require.True(t, opt.Empty(), "opt should be empty")

	// initializes iterator
	err = iter.Next(ctx)
	require.NoError(t, err)

	for i := 0; !iter.Finished(); i++ {
		opt := iter.Value()
		t.Logf("%v", opt)
		require.False(t, opt.Empty())

		want := testData[i]
		require.Equal(t, want, string(iter.CurrentKey()))
		require.Equal(t, want, string(opt.GetValue()))

		err := iter.Next(ctx)
		require.NoError(t, err)
	}

	require.Nil(t, iter.CurrentKey())
	opt = iter.Value()
	require.True(t, opt.Empty())

	err = iter.Close()
	require.NoError(t, err)
}

// RunSuite runs a test harness for an implementation of store.KVStore
// producer is a function which MUST return a new instance of a store.KVStore for each call
func RunSuite(t *testing.T, producer func() store.KVStore) {
	tests := []func(*testing.T, store.KVStore){
		test_DeleteKeyNotSet_ReturnsRemovedFalse,
		test_DeleteNilKey_Errors,
		test_GetDocumentNotRegistered_ReturnsEmptyOption,
		test_GetWithNilKey_Errors,
		test_HasNilKey_Errors,
		test_Has_FalseWhenNotSet,
		test_Has_TrueWhenSet,
		test_Iterate_ReturnsIteratorOverAllItems,
		test_SetNewValue_ReturnsCreatedTrue,
		test_SetNilKey_Errors,
		test_SetThenDelete_KeyNoLongerInStore,
		test_SettingThenGet_ReturnsDoc,
		test_UpdateKey_ReturnsCreatedFalse,
		test_SetNilKV_ReturnsNilKeyErr,
		test_SetNilValue_Ok,
	}

	for _, test := range tests {
		funcPtr := reflect.ValueOf(test).Pointer()
		fullName := runtime.FuncForPC(funcPtr).Name()
		parts := strings.Split(fullName, ".")
		name := parts[len(parts)-1]
		kv := producer()
		t.Run(name, func(t *testing.T) {
			test(t, kv)
		})
	}
}
