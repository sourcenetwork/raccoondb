package iterator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_MapIter_AppliesMap(t *testing.T) {
	vals := []int{0, 1, 2, 3}
	iter := FromSlice(vals)

	iter = Map(iter, func(x int) int {
		return x + 1
	})

	got, errs := Consume(iter)

	require.Empty(t, errs)
	want := []int{1, 2, 3, 4}
	require.Equal(t, want, got)
}

/*
func Test_MapFailableIter_AppliesMap(t *testing.T) {
	vals := []int{0, 1, 2, 3}
	iter := FromSlice(vals)

	iter = MapFailable(iter, func(x int) (int, error) {
		if x%2 == 0 {
			return 0, errors.New("even")
		}
		return x, nil
	})

	for i:=0; !iter.Finished(); i++ {
		opt := iter.Value()
		if i%2 == 0 {
			require.True(t, opt.Empty())

		}
	}
	got, errs := Consume(iter)

	require.Empty(t, errs)
	want := []int{1, 2, 3, 4}
	require.Equal(t, want, got)

}

*/
