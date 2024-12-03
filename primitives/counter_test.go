package primitives

import (
	"context"
	"testing"

	"github.com/sourcenetwork/raccoondb/v2/store/corekv"
	"github.com/stretchr/testify/require"
)

var testCounter []byte = []byte("counter")

func Test_CounterStore_GetSet(t *testing.T) {
	ctx := context.TODO()
	kv := corekv.NewMemKV()
	counter := NewCounterStore(kv)

	val, err := counter.Get(ctx, testCounter)
	require.NoError(t, err)
	require.Equal(t, uint64(0), val)

	val, err = counter.Increment(ctx, testCounter)
	require.NoError(t, err)
	require.Equal(t, uint64(1), val)

	val, err = counter.Get(ctx, testCounter)
	require.NoError(t, err)
	require.Equal(t, uint64(1), val)
}

func Test_CounterStore_GetNext_ReturnsNextNumber(t *testing.T) {
	ctx := context.TODO()
	kv := corekv.NewMemKV()
	counter := NewCounterStore(kv)

	_, err := counter.Increment(ctx, testCounter)
	require.NoError(t, err)
	_, err = counter.Increment(ctx, testCounter)
	require.NoError(t, err)

	val, err := counter.GetNext(ctx, testCounter)
	require.NoError(t, err)
	require.Equal(t, uint64(3), val)

	val, err = counter.Get(ctx, testCounter)
	require.NoError(t, err)
	require.Equal(t, uint64(2), val)
}
func Test_CounterStore_GetNext_Returns1IfUninitialized(t *testing.T) {
	ctx := context.TODO()
	kv := corekv.NewMemKV()
	counter := NewCounterStore(kv)

	val, err := counter.GetNext(ctx, testCounter)
	require.NoError(t, err)
	require.Equal(t, uint64(1), val)
}

func Test_CounterStore_Get_Returns0IfUninitialized(t *testing.T) {
	ctx := context.TODO()
	kv := corekv.NewMemKV()
	counter := NewCounterStore(kv)

	val, err := counter.Get(ctx, testCounter)
	require.NoError(t, err)
	require.Equal(t, uint64(0), val)
}

func Test_CounterStore_GetSetMulti(t *testing.T) {
	ctx := context.TODO()
	kv := corekv.NewMemKV()
	counter := NewCounterStore(kv)

	key2 := []byte("othercounter")

	_, err := counter.Increment(ctx, testCounter)
	require.NoError(t, err)
	_, err = counter.Increment(ctx, testCounter)
	require.NoError(t, err)
	_, err = counter.Increment(ctx, key2)
	require.NoError(t, err)

	val, err := counter.Get(ctx, testCounter)
	require.NoError(t, err)
	require.Equal(t, uint64(2), val)

	val, err = counter.Get(ctx, key2)
	require.NoError(t, err)
	require.Equal(t, uint64(1), val)
}

func Test_CounterStore_Has_FalseIfNotInit(t *testing.T) {
	ctx := context.TODO()
	kv := corekv.NewMemKV()
	counter := NewCounterStore(kv)

	has, err := counter.Has(ctx, testCounter)
	require.NoError(t, err)
	require.False(t, has)
}

func Test_CounterStore_Decrement_Reduces(t *testing.T) {
	ctx := context.TODO()
	kv := corekv.NewMemKV()
	counter := NewCounterStore(kv)

	_, err := counter.Increment(ctx, testCounter)
	require.NoError(t, err)
	_, err = counter.Increment(ctx, testCounter)
	require.NoError(t, err)

	new, err := counter.Decrement(ctx, testCounter)
	require.NoError(t, err)
	require.Equal(t, uint64(1), new)

	stored, err := counter.Get(ctx, testCounter)
	require.NoError(t, err)
	require.Equal(t, uint64(1), stored)
}

func Test_CounterStore_DeleteCounter(t *testing.T) {
	ctx := context.TODO()
	kv := corekv.NewMemKV()
	counter := NewCounterStore(kv)

	_, err := counter.Increment(ctx, testCounter)
	require.NoError(t, err)

	removed, err := counter.Delete(ctx, testCounter)
	require.NoError(t, err)
	require.True(t, bool(removed))

	has, err := counter.Has(ctx, testCounter)
	require.NoError(t, err)
	require.False(t, has)
}
