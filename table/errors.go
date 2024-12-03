package table

import (
	"fmt"

	"github.com/sourcenetwork/raccoondb/v2/errors"
)

var ErrIndexedObjectStore = errors.New("indexed object store")
var ErrIndexExists = errors.New("index already defined")
var ErrObjectIndex = errors.New("ObjectIndex")

func newIndexErr(method string, msg string, err error) error {
	return fmt.Errorf("%w: %v: %v: %w", ErrObjectIndex, method, msg, err)
}

func newTableErr(method string, msg string, err error) error {
	return fmt.Errorf("%w: %v: %v: %w", ErrIndexedObjectStore, method, msg, err)
}
