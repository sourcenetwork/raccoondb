package iterator

import (
	"context"

	"github.com/sourcenetwork/raccoondb/v2/types"
)

// Iterator models a stateful traversing through some sequence of elements
// indexed by a byte sequence key
//
// A new iterator should start "out of bound", meanign it does not perform
// any IO until Next() is called for the first time.
// The iterator ends once it steps out of bound.
//
// Expected usage:
//
// defer iter.Close()
// err := iter.Next()
// if err != nil { return err }
//
//	for !iter.Finished() {
//	  _ = iter.Value()
//	  err := iter.Next();
//	  if err != nil { return err }
//	}
//
//		or:
//
// defer iter.Close()
//
//	for {
//	  err := iter.Next()
//	  if err != nil {
//	    return err
//	  }
//	  if iter.Finished() {
//	    break
//	  }
//	  _ = iter.Value()
//	}
type Iterator[T any] interface {
	// Next steps the iterator to its next value.
	// It may return an error if it could not produce a value.
	// The Iterator will yield values until Finished == true,
	// once an Iteartor is finished, Next MUST be a Noop and return no error.
	//
	// Note: An error does not necessarily mean that the Iterator is Finished.
	Next(ctx context.Context) error

	// Value returns the current value in the Iterator
	// Should only return None if Next returned an error or if the Iterator is Finished
	Value() types.Option[T]

	// Finished indicates whether the Iterator scanned through all possible keys
	Finished() bool

	// Close frees up resources taken by the Iterator
	Close() error

	// GetParams returns the control options specifying the iterator
	GetParams() IteratorOpt

	// CurrentKey returns the key of the current element
	// If Finished is true, return nil
	CurrentKey() []byte
}

// BytesIterator is a type alias for an iterator which returns a sequence of byte slices
type BytesIterator Iterator[[]byte]
