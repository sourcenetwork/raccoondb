package iterator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SeekKeyPrefix_StepsUntilPrefix(t *testing.T) {
	items := map[string]string{
		"001": "001",
		"111": "111",
		"222": "222",
	}
	iter := IterFromStringKeyMap(items)

	found, err := SeekKeyPrefix(context.TODO(), iter, []byte("1"))

	require.Equal(t, "111", string(iter.CurrentKey()))
	require.NoError(t, err)
	require.True(t, found)
}

func Test_SeekKeyPrefix_ConsumesIterIfKeysSmallerThanPrefix(t *testing.T) {
	items := map[string]string{
		"001": "001",
		"002": "002",
		"003": "003",
	}
	iter := IterFromStringKeyMap(items)

	found, err := SeekKeyPrefix(context.TODO(), iter, []byte("1"))

	// Then seek consumes iterator because all keys are less than 1
	require.False(t, found)
	require.NoError(t, err)
	require.True(t, iter.Finished())
}

func Test_SeekKeyPrefix_StopsOnceKeyIsLargerThanPrefix(t *testing.T) {
	items := map[string]string{
		"001": "001",
		"002": "002",
		"003": "003",
		"222": "222",
		"223": "223",
	}
	iter := IterFromStringKeyMap(items)

	found, err := SeekKeyPrefix(context.TODO(), iter, []byte("1"))

	require.False(t, found)
	require.NoError(t, err)
	// Then iterator stopped at 222, because 2 > 1
	require.Equal(t, "222", string(iter.CurrentKey()))
	require.False(t, iter.Finished())
}

func Test_SeekKeyPrefix_IfPrefixMatchesKey_StopAtKey(t *testing.T) {
	items := map[string]string{
		"0": "0",
		"1": "1",
		"2": "2",
	}
	iter := IterFromStringKeyMap(items)

	found, err := SeekKeyPrefix(context.TODO(), iter, []byte("1"))

	require.True(t, found)
	require.NoError(t, err)
	require.Equal(t, "1", string(iter.CurrentKey()))
	require.False(t, iter.Finished())
}
