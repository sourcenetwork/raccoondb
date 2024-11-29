package iterator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IterFromSlice_IteratesOverAllValues(t *testing.T) {
	vals := []int{
		1,
		2,
		3,
		4,
		5,
	}

	iter := IterFromSlice(vals)

	got, errs := Consume(context.TODO(), iter)
	require.Empty(t, errs)
	require.Equal(t, vals, got)
}
