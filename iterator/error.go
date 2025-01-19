package iterator

import "strings"

// IterItemError models an IterItemError which happened at a given key
type IterItemError struct {
	Key []byte
	Err error
}

// IterationError bundles IterItemError
type IterationError struct {
	Errors []IterItemError
}

func (e *IterationError) Error() string {
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
