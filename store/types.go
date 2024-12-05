// package store contains type definitions for Raccoon's KV Store
package store

import (
	"context"

	"github.com/sourcenetwork/raccoondb/v2/iterator"
	"github.com/sourcenetwork/raccoondb/v2/types"
)

// KeyRemoved is a flag indicating whether a Key was removed during a Delete call
type KeyRemoved bool

// KeyCreated is a flag indicating whether a Key was created during a Set call
type KeyCreated bool

// StoreIterator models an Iterator aware
// of the Iteration params
type StoreIterator[T any] interface {
	iterator.Iterator[T]
	// GetParams returns the control options specifying the iterator
	GetParams() IteratorOpt
}

// Iterable models a store which can be iterated over
type Iterable[T any] interface {
	// Iterate returns an iterator which walks through the stored pairs
	Iterate(ctx context.Context, opt IteratorOpt) (StoreIterator[T], error)
}

// ReadStore is a subset of KVStore which contains read only methods
type ReadStore interface {
	Iterable[[]byte]

	// Get fetches the value associated to the given key
	// Return Option None if the key was not found
	Get(ctx context.Context, key []byte) (types.Option[[]byte], error)

	// Has checks whether key has an associated value
	Has(ctx context.Context, key []byte) (bool, error)
}

// KVStore models a Key-Value store
type KVStore interface {
	ReadStore

	// Set stores a key-value pair in the store
	// Return RecordCreated true if a new node / entry was created in the underlying store
	Set(ctx context.Context, key, value []byte) (KeyCreated, error)

	// Delete removes an entry for the key-value store.
	// Return RecordRemoved false if the record was not found
	Delete(ctx context.Context, key []byte) (KeyRemoved, error)
}
