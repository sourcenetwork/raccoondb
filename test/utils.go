package test

import (
	"github.com/sourcenetwork/raccoondb/v2/store"
	"github.com/sourcenetwork/raccoondb/v2/store/corekv"
)

func NewTestKV() store.KVStore {
	return corekv.NewMemKV()

}
