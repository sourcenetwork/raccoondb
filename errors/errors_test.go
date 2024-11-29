package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ErrorsAs_ReturnsTrueForRaccoonError(t *testing.T) {
	err := New("test error")

	var rerr *RaccoonError
	ok := errors.As(err, &rerr)

	require.True(t, ok)
}
