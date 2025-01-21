package iterator

import (
	"context"
)

// Iterator models a stateful traversing through some sequence of elements
// indexed by a byte sequence key
//
// A new iterator will start ready to produce the first value.
// The iterator ends once it steps out of bound.
//
// Example usage:
//
// defer iter.Close()
//
//	for !iter.Finished() {
//	  _ = iter.Value()
//	  err := iter.Next();
//	  if err != nil { return err }
//	}
type Iterator[T any] interface {
	// Next steps the iterator to its next value.
	// It may return an error if it could not step the iterator.
	// The Iterator will yield values until Finished == true,
	// once an Iteartor is finished, Next MUST be a Noop and return no error.
	//
	// Note: An error does not necessarily mean that the Iterator is Finished.
	Next(ctx context.Context) error

	// Value returns the current value in the Iterator or an error
	Value() (T, error)

	// Finished indicates whether the Iterator scanned through all possible keys
	Finished() bool

	// Close frees up resources taken by the Iterator
	Close() error

	// CurrentKey returns the key of the current element
	// If Finished is true, return nil
	CurrentKey() []byte
}

// BytesIterator is a type alias for an iterator which returns a sequence of byte slices
type BytesIterator Iterator[[]byte]

// KeyValue models a key value pair
type KeyValue[T any] struct {
	Key   []byte
	Value T
}
