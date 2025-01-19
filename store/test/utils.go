package test

import (
	"context"
	"testing"

	"github.com/sourcenetwork/raccoondb/v2/iterator"
	"github.com/sourcenetwork/raccoondb/v2/store"
)

func DumpStore(t *testing.T, kv store.KVStore) {
	iter, _ := kv.Iterate(context.TODO(), store.NewOpenIterator())
	pairs, err := iterator.ConsumePairs(context.TODO(), iter)
	if err != nil {
		panic(err)
	}
	for _, pair := range pairs {
		t.Logf("\tkey: %v\t value: %v", string(pair.Key), pair.Value)
	}
}
