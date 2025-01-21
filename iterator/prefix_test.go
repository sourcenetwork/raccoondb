package iterator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Prefix_WhithoutPrefix_ReturnsNoItems(t *testing.T) {
	items := map[string]string{
		"abc":  "abc",
		"def":  "def",
		"test": "test",
	}
	iter := IterFromStringKeyMap(items)
	iter, err := NewPrefixIterator(context.TODO(), []byte("_"), iter, false)
	require.NoError(t, err)

	vals, errs := Consume(context.TODO(), iter)
	require.NoError(t, errs)
	require.Len(t, vals, 0)
}

func Test_Prefix_WithPrefixedItems_ReturnItemsUntilEnd(t *testing.T) {
	items := map[string]string{
		"_abc":  "abc",
		"_def":  "def",
		"_test": "test",
	}
	iter := IterFromStringKeyMap(items)
	iter, err := NewPrefixIterator(context.TODO(), []byte("_"), iter, false)
	require.NoError(t, err)

	vals, err := Consume(context.TODO(), iter)

	require.NoError(t, err)
	want := []string{
		"abc",
		"def",
		"test",
	}
	require.Equal(t, want, vals)
}

func Test_Prefix_WithItemsBeforePrefix_IterIsCreatedAfterItemsThatDoNotContainPrefix(t *testing.T) {
	items := map[string]string{
		"001": "001",
		"111": "111",
		"112": "112",
	}
	iter := IterFromStringKeyMap(items)
	iter, err := NewPrefixIterator(context.TODO(), []byte("1"), iter, false)
	require.NoError(t, err)

	val, err := iter.Value()
	require.NoError(t, err)
	require.Equal(t, "111", val)
}

func Test_Prefix_WithItemsAfterPrefix_ReturnsNoItemsWithoutPrefix(t *testing.T) {
	items := map[string]string{
		"001": "001",
		"111": "111",
		"112": "112",
		"222": "222",
		"333": "333",
	}
	iter := IterFromStringKeyMap(items)
	iter, err := NewPrefixIterator(context.TODO(), []byte("1"), iter, false)
	require.NoError(t, err)

	vals, errs := Consume(context.TODO(), iter)
	require.Empty(t, errs)
	want := []string{"111", "112"}
	require.Equal(t, want, vals)
}
