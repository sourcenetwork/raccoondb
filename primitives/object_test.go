package primitives

import (
	"testing"

	"github.com/sourcenetwork/raccoondb/v2/marshal"
	"github.com/sourcenetwork/raccoondb/v2/store"
	"github.com/sourcenetwork/raccoondb/v2/store/test"
	testutil "github.com/sourcenetwork/raccoondb/v2/test"
)

func Test_KeyObjectStore_Suite(t *testing.T) {
	m := &marshal.BytesMarshaler{}
	producer := func() store.KVStore {
		kv := testutil.NewTestKV()
		okv := NewKeyObjectStore(kv, m)
		return &okv
	}
	test.RunSuite(t, producer)
}
