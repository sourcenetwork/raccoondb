package table

import (
	"fmt"

	"github.com/sourcenetwork/raccoondb/v2/errors"
)

// ErrTable is the top level error for Table errors
var ErrTable = errors.New("table")

// ErrIndexExists is thrown when an index with conflicting name is created to the same table
var ErrIndexExists = errors.New("index already defined")

// ErrTableIndex is thrown when an error happens reading / writing to and from a Table index
var ErrTableIndex = errors.New("table index")

func newIndexErr(method string, msg string, err error) error {
	return fmt.Errorf("%w: %v: %v: %w", ErrTableIndex, method, msg, err)
}

func newTableErr(method string, msg string, err error) error {
	return fmt.Errorf("%w: %v: %v: %w", ErrTable, method, msg, err)
}
