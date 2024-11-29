package test

import (
	"context"
	"testing"

	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/store"
)

func DumpStore(t *testing.T, kv store.KVStore) {
	iter, _ := kv.Iterate(context.TODO(), iterator.NewOpenIterator())
	pairs := iterator.ConsumePairs(context.TODO(), iter)
	for _, pair := range pairs {
		t.Logf("\tkey: %v\t value: %v", string(pair.Key), pair.Value)
	}
}
