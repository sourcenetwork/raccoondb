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

// ReadStore is a subset of KVStore which contains read only methods
type ReadStore interface {
	// Get fetches the value associated to the given key
	// Return Option None if the key was not found
	Get(ctx context.Context, key []byte) (types.Option[[]byte], error)

	// Has checks whether key has an associated value
	Has(ctx context.Context, key []byte) (bool, error)

	// Iterate returns an iterator which walks through the stored pairs
	Iterate(ctx context.Context, opt iterator.IteratorOpt) (iterator.BytesIterator, error)
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
