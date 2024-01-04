package raccoon

import (
    "cosmossdk.io/store/mem"
    "cosmossdk.io/store/types"
    tmdb "github.com/tendermint/tm-db"
)

var _ KVStore = (*tmdbWrapper)(nil)

func NewMemKV() KVStore {
    return KvFromCosmosKv(mem.NewStore())
}

func NewLevelDB(path, file string) (KVStore, error) {
    db, err := tmdb.NewGoLevelDB(file, path)
    if err != nil {
        return nil, err
    }
    wrapper := tmdbWrapper{
        db: db,
    }
    return &wrapper, nil
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

type tmdbWrapper struct {
    db tmdb.DB
}

func (s *tmdbWrapper) Get(key []byte) ([]byte, error) {
    return s.db.Get(key)
}

func (s *tmdbWrapper) Has(key []byte) (bool, error) {
    return s.db.Has(key)
}

func (s *tmdbWrapper) Set(key []byte, val []byte) error {
   return s.db.Set(key, val)
}

func (s *tmdbWrapper) Delete(key []byte) error {
    return s.db.Delete(key)
}

func (s *tmdbWrapper) Iterator(start, end []byte) Iterator {
    iter, _ := s.db.Iterator(start, end)
    return iter
}
