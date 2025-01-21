package iterator

import (
	"context"
	"testing"

	"github.com/sourcenetwork/raccoondb/v2/errors"
	"github.com/stretchr/testify/require"
)

func Test_MapIter_AppliesMap(t *testing.T) {
	vals := []int{0, 1, 2, 3}
	iter := IterFromSlice(vals)

	iter = Map(iter, func(x int) int {
		return x + 1
	})

	got, errs := Consume(context.TODO(), iter)

	require.Empty(t, errs)
	want := []int{1, 2, 3, 4}
	require.Equal(t, want, got)
}

func Test_MapFailableIter_AppliesMap(t *testing.T) {
	vals := []int{0, 1, 2, 3}
	iter := IterFromSlice(vals)

	iter = MapFailable(iter, func(x int) (int, error) {
		if x%2 == 0 {
			return 0, errors.New("even")
		}
		return x, nil
	})

	ctx := context.TODO()

	for i := 0; !iter.Finished(); i++ {
		val, err := iter.Value()
		if i%2 == 0 {
			require.Error(t, err, "index %v", i)
			require.Equal(t, 0, val)
		} else {
			require.NoError(t, err)
			require.Equal(t, i, val)
		}

		err = iter.Next(ctx)
		require.NoError(t, err)

		if iter.Finished() {
			break
		}
	}
}
