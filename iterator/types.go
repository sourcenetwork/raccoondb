package iterator

import "github.com/sourcenetwork/raccoondb/types"

// Iterator models a stateful traversing through some sequence of elements
// indexed by a byte sequence key
//
// Expected usage:
//
//	defer iter.Close()
//	for iter.Finished() {
//		val := iter.Value()
//		err := iter.Next();
//		if err != nil {
//			return err
//		}
//	}
//
// Meaning that the iterator starts ready to supply a value,
// next moves it forward for as long as values are valid,
// and the final next call moves it to out of bounds where it is no longer valid and finished turns true
type Iterator[T any] interface {
	// Next steps the iterator to its next value.
	// It may error if for some reason the value cannot be produced.
	// An error does not necessarily mean that the iterator is Finished.
	// If the Iterator is Finished, Next is a Noop and returns no error
	Next() error

	// Value returns the current value in the Iterator
	// Should only return None if Next returned an error
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
