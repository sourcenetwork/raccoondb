package corekv

import (
	"context"
	"testing"

	"github.com/sourcenetwork/corekv/memory"
	"github.com/sourcenetwork/raccoondb/store"
	"github.com/sourcenetwork/raccoondb/store/test"
)

func Test_corekv_RunSuite(t *testing.T) {
	producer := func() store.KVStore {
		corekv := memory.NewDatastore(context.TODO())
		kv := WrapCoreKV(corekv)
		return kv
	}
	test.RunSuite(t, producer)
}
