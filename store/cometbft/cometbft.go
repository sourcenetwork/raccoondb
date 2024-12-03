// package cometbft defines an adaptor which wraps cometbft DB into a Raccoon KV Store
package cometbft

import (
	cmdb "github.com/cometbft/cometbft-db"

	"github.com/sourcenetwork/raccoondb/v2/store"
)

// NewMemKV returns a cometbft memory store adapted to Racoon's KVStore
func NewMemKV() store.KVStore {
	return KVFromCometDB(cmdb.NewMemDB())
}
