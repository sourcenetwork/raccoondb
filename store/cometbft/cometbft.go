package cometbft

import (
	cmdb "github.com/cometbft/cometbft-db"

	"github.com/sourcenetwork/raccoondb/v2/store"
)

func NewMemKV() store.KVStore {
	return KVFromCometDB(cmdb.NewMemDB())
}
