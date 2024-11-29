package test

import (
	"context"
	"testing"

	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/store"
)

func DumpStore(t *testing.T, kv store.KVStore) {
	iter, _ := kv.Iterate(context.TODO(), iterator.NewOpenIterator())
	keys := iterator.ConsumeKeys(context.TODO(), iter)
	for _, key := range keys {
		t.Logf("\tkey: %v\n", string(key))
	}
}
