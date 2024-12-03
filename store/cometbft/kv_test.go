package cometbft

import (
	"testing"

	cmdb "github.com/cometbft/cometbft-db"

	"github.com/sourcenetwork/raccoondb/v2/store"
	"github.com/sourcenetwork/raccoondb/v2/store/test"
)

func Test_Comet_RunSuite(t *testing.T) {
	factory := func() store.KVStore {
		db := cmdb.NewMemDB()
		return KVFromCometDB(db)
	}
	test.RunSuite(t, factory)
}
