package iterator

import "github.com/sourcenetwork/raccoondb/types"

// Iterator models a stateful traversing through some sequence of elements
// indexed by a byte sequence key
type Iterator[T any] interface {
	// Next steps the iterator to its next value.
	// It may error if for some reason the value cannot be produced.
	// An error does not necessarily mean that the iterator is Finished.
	// If the Iterator is Finished, Next is a Noop and returns no error
	Next() error

	// Value returns the current value in the Iterator
	Value() types.Option[T]

	// Finished indicates whether the Iterator scanned through all possible keys
	Finished() bool

	// Close frees up resources taken by the Iterator
	Close() error

	// GetParams returns the control options specifying the iterator
	GetParams() IteratorOpt

	// CurrentKey returns the key of the current element
	CurrentKey() []byte
}

// BytesIterator is a type alias for an iterator which returns a sequence of byte slices
type BytesIterator Iterator[[]byte]

// IteratorOpt configures the behavior of an Iterator
// TODO improve this UX
type IteratorOpt struct {
	// Start represents the lower bound of iteration
	// If nil will start at the smallest element
	Start []byte

	// End represents the uper bound of iteration
	// If nil will end at the largest element
	End []byte

	// Prefix does a prefix iteration on the store
	Prefix []byte

	// Reverse iterates the store backwards
	Reverse bool
}
