package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_errorsAs_ReturnsTrueForRaccoonError(t *testing.T) {
	err := New("test error")

	raccoonErr := RaccoonError{}
	ok := errors.As(err, &raccoonErr)

	require.True(t, ok)
}
