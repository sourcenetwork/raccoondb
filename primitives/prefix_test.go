package primitives

import (
	"testing"

	"github.com/sourcenetwork/raccoondb/store"
	"github.com/sourcenetwork/raccoondb/store/corekv"
	"github.com/sourcenetwork/raccoondb/store/test"
)

func Test_PrefixKV_Suite(t *testing.T) {
	prefix := []byte("prefix")
	factory := func() store.KVStore {
		kv := corekv.NewMemKV()
		return NewPrefixedKV(kv, prefix)
	}
	test.RunSuite(t, factory)
}
