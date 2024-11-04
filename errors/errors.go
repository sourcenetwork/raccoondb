package errors

import "errors"

var (
	ErrKV      = errors.New("kv store error")
	ErrMarshal = errors.New("marshaling error")
)
