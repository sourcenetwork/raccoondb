package corekv

import (
	"fmt"

	"github.com/sourcenetwork/raccoondb/v2/errors"
)

// ErrCoreKV is a base error for adapted corekv errors
var ErrCoreKV = errors.New("corekv")

// wrapErr wraps an error with ErrCoreKV
func wrapErr(err error) error {
	return fmt.Errorf("%w: %w", ErrCoreKV, err)
}
