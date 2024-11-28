// package errors defines the top level error type for raccoondb
package errors

var _ error = (*RaccoonError)(nil)

// RaccoonError is a simple error type used by raccoon
// All errors returned by raccoon can be converted to a RaccoonError
// using errors.As
type RaccoonError struct {
	msg string
}

// Error implements the error interface
func (e *RaccoonError) Error() string {
	return e.msg
}

// New returns a new RaccoonError
func New(msg string) error {
	return &RaccoonError{
		msg: msg,
	}
}

// As returns true and sets e.msg to err.msg
// if err is an instances RaccoonError
func (e *RaccoonError) As(err error) bool {
	cast, ok := err.(*RaccoonError)
	if ok {
		e.msg = cast.msg
	}
	return ok
}
