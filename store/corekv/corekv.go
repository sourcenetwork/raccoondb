// package corekv defines an adaptor which wraps corekv KVstore's into Raccoon KV Stores
package corekv

import (
	"context"

	"github.com/sourcenetwork/corekv/memory"
	"github.com/sourcenetwork/raccoondb/v2/store"
)

// NewMemKV returns a corekv memory store adapted to Racoon's KVStore
func NewMemKV() store.KVStore {
	mem := memory.NewDatastore(context.TODO())
	return WrapCoreKV(mem)
}
