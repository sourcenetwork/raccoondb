package iterator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Filter_FiltersItems(t *testing.T) {
	vals := []int{1, 2, 3, 4, 5, 6}
	iter := IterFromSlice(vals)

	// Filter for even numbers
	iter = Filter(iter, func(x int) bool {
		return x%2 == 0
	})

	got, err := Consume(context.TODO(), iter)
	require.NoError(t, err)
	want := []int{2, 4, 6}
	require.Equal(t, want, got)
}

func Test_Filter_NoMatchingItems_ReturnsEmpty(t *testing.T) {
	vals := []int{1, 3, 5, 7}
	iter := IterFromSlice(vals)

	// Filter for even numbers (none exist)
	iter = Filter(iter, func(x int) bool {
		return x%2 == 0
	})

	got, err := Consume(context.TODO(), iter)
	require.NoError(t, err)
	require.Empty(t, got)
}

func Test_Filter_EmptyIterator_ReturnsEmpty(t *testing.T) {
	iter := NewEmptyIterator[int]()

	iter = Filter(iter, func(x int) bool {
		return true
	})

	got, err := Consume(context.TODO(), iter)
	require.NoError(t, err)
	require.Empty(t, got)
}

func Test_Filter_CurrentKey_DoesNotCauseInfiniteRecursion(t *testing.T) {
	items := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	iter := IterFromStringKeyMap(items)

	iter = Filter(iter, func(x int) bool {
		return x%2 == 1 // odd numbers
	})

	// Iterator is now initialized on creation, positioned on first match
	// This should not cause infinite recursion
	key := iter.CurrentKey()
	require.NotNil(t, key)
}

func Test_Filter_ConsumePairs_ReturnsCorrectKeyValuePairs(t *testing.T) {
	items := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}
	iter := IterFromStringKeyMap(items)

	// Filter for even numbers
	iter = Filter(iter, func(x int) bool {
		return x%2 == 0
	})

	ctx := context.TODO()
	pairs, err := ConsumePairs(ctx, iter)
	require.NoError(t, err)
	require.Len(t, pairs, 2)

	// Verify keys match their values
	for _, pair := range pairs {
		key := string(pair.Key)
		val := pair.Value
		if key == "b" {
			require.Equal(t, 2, val)
		} else if key == "d" {
			require.Equal(t, 4, val)
		} else {
			t.Errorf("unexpected key: %s", key)
		}
	}
}

func Test_Filter_CurrentKey_ReturnsCorrectKeys(t *testing.T) {
	items := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}
	iter := IterFromStringKeyMap(items)

	// Filter for odd numbers
	iter = Filter(iter, func(x int) bool {
		return x%2 == 1
	})

	ctx := context.TODO()
	var keys []string

	// Iterator is initialized on creation, already positioned on first match
	for !iter.Finished() {
		key := iter.CurrentKey()
		keys = append(keys, string(key))
		err := iter.Next(ctx)
		require.NoError(t, err)
	}

	// Should have keys for odd values (1 and 3)
	require.Len(t, keys, 2)
	require.Contains(t, keys, "a") // value 1
	require.Contains(t, keys, "c") // value 3
}
