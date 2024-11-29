package iterator

import "strings"

// IterationError models an IterationError which happened at a given key
type IterationError struct {
	Key []byte
	Err error
}

// SeekError bundles IterationErrors
type SeekError struct {
	Errors []IterationError
}

func (e *SeekError) Error() string {
	builder := strings.Builder{}
	for _, err := range e.Errors {
		builder.WriteString("key ")
		builder.Write(err.Key)
		builder.WriteString(" : ")
		builder.WriteString(err.Err.Error())
		builder.WriteRune('\n')
	}
	return builder.String()
}
