package iterator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TakeNIterator_YieldsSetNumberOfElements(t *testing.T) {
	items := map[string]string{
		"001": "001",
		"111": "111",
		"112": "112",
		"222": "222",
		"333": "333",
	}
	iter := IterFromStringKeyMap(items)

	iter, err := NewTakeNIterator(iter, 2)
	require.NoError(t, err)

	vals, err := Consume(context.TODO(), iter)
	require.NoError(t, err)
	want := []string{
		"001", "111",
	}
	require.Equal(t, want, vals)
}

func Test_TakeNIterator_ErrorsIfCount0(t *testing.T) {
	items := map[string]string{
		"001": "001",
		"111": "111",
		"112": "112",
		"222": "222",
		"333": "333",
	}
	iter := IterFromStringKeyMap(items)

	iter, err := NewTakeNIterator(iter, 0)
	require.Nil(t, iter)
	require.ErrorIs(t, err, ErrZeroCount)
}

func Test_TakeNIterator_Returns1ElemWithCount1(t *testing.T) {
	items := map[string]string{
		"001": "001",
		"111": "111",
		"112": "112",
		"222": "222",
		"333": "333",
	}
	iter := IterFromStringKeyMap(items)

	iter, err := NewTakeNIterator(iter, 1)
	require.NoError(t, err)

	vals, err := Consume(context.TODO(), iter)
	require.NoError(t, err)
	want := []string{
		"001",
	}
	require.Equal(t, want, vals)
}
