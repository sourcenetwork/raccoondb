package primitives

import (
	"context"

	"github.com/sourcenetwork/raccoondb/v2/store"
)

// CountedStore models a generic storage system which keeps track of the amount of entries it has
type CountedStore interface {
	// GetCount returns the number of entires the KVStore has
	GetCount(context.Context) (uint64, error)
}

// CountedKVStore models a KVStore which keeps track of the amount of values it stores
type CountedKVStore interface {
	store.KVStore
	CountedStore
}
