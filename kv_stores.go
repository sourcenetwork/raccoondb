package raccoon

import (
    "github.com/cosmos/cosmos-sdk/store/mem"
    "github.com/cosmos/cosmos-sdk/store/types"
)


func NewMemKV() KVStore {
    return KvFromCosmosKv(mem.NewStore())
}

func KvFromCosmosKv(store types.KVStore) KVStore {
    return &cosmosKvWrapper {
        store: store,
    }
}

type cosmosKvWrapper struct {
    store types.KVStore
}

func (s *cosmosKvWrapper) Get(key []byte) ([]byte, error) {
    return s.store.Get(key), nil
}

func (s *cosmosKvWrapper) Has(key []byte) (bool, error) {
    return s.store.Has(key), nil
}

func (s *cosmosKvWrapper) Set(key []byte, val []byte) error {
   s.store.Set(key, val) 
   return nil
}

func (s *cosmosKvWrapper) Delete(key []byte) error {
    s.store.Delete(key)
    return nil
}

func (s *cosmosKvWrapper) Iterator(start, end []byte) Iterator {
    return s.store.Iterator(start, end)
}
