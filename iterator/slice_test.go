package iterator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FromSlice_IteratesOverAllValues(t *testing.T) {
	vals := []int{
		1,
		2,
		3,
		4,
		5,
	}

	iter := FromSlice(vals)

	got, errs := Consume(iter)
	require.Empty(t, errs)
	require.Equal(t, vals, got)
}
