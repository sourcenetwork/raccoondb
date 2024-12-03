package cometbft

import (
	"fmt"

	"github.com/sourcenetwork/raccoondb/v2/errors"
)

var ErrCometDB = errors.New("cometbft-db")

// wrapErr wraps an error with ErrCometDb
func wrapErr(err error) error {
	return fmt.Errorf("%w: %w", ErrCometDB, err)
}
