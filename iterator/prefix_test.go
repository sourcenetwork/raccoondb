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
	iter = NewPrefixIterator([]byte("_"), iter, false)

	vals, errs := Consume(context.TODO(), iter)

	require.Len(t, vals, 0, "vals got %v", vals)
	require.Len(t, errs, 0, "err got %v", vals)
}

func Test_Prefix_WithPrefixedItems_ReturnItemsUntilEnd(t *testing.T) {
	items := map[string]string{
		"_abc":  "abc",
		"_def":  "def",
		"_test": "test",
	}
	iter := IterFromStringKeyMap(items)
	iter = NewPrefixIterator([]byte("_"), iter, false)

	vals, errs := Consume(context.TODO(), iter)

	want := []string{
		"abc",
		"def",
		"test",
	}
	require.Equal(t, want, vals)
	require.Len(t, errs, 0, "err got %v", vals)
}

func Test_Prefix_WithItemsBeforePrefix_NextSkipsThatDoNotContainPrefix(t *testing.T) {
	items := map[string]string{
		"001": "001",
		"111": "111",
		"112": "112",
	}
	iter := IterFromStringKeyMap(items)
	iter = NewPrefixIterator([]byte("1"), iter, false)

	err := iter.Next(context.TODO())
	require.NoError(t, err)

	opt := iter.Value()
	require.False(t, opt.Empty())
	require.Equal(t, "111", opt.GetValue())
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
	iter = NewPrefixIterator([]byte("1"), iter, false)

	vals, errs := Consume(context.TODO(), iter)
	require.Empty(t, errs)
	want := []string{"111", "112"}
	require.Equal(t, want, vals)
}
