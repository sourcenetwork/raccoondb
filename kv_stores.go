package raccoon

import (
	storetypes "cosmossdk.io/store/types"

	"cosmossdk.io/store/mem"
	dbm "github.com/cosmos/cosmos-db"
)

var _ KVStore = (*dbmWrapper)(nil)

func NewMemKV() KVStore {
	return KvFromCosmosKv(mem.NewStore())
}

func NewLevelDB(path, file string) (KVStore, error) {
	db, err := dbm.NewGoLevelDB(file, path, nil)
	if err != nil {
		return nil, err
	}
	wrapper := dbmWrapper{
		db: db,
	}
	return &wrapper, nil
}

func KvFromCosmosKv(store storetypes.KVStore) KVStore {
	return &cosmosKvWrapper{
		store: store,
	}
}

type cosmosKvWrapper struct {
	store storetypes.KVStore
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

type dbmWrapper struct {
	db dbm.DB
}

func (s *dbmWrapper) Get(key []byte) ([]byte, error) {
	return s.db.Get(key)
}

func (s *dbmWrapper) Has(key []byte) (bool, error) {
	return s.db.Has(key)
}

func (s *dbmWrapper) Set(key []byte, val []byte) error {
	return s.db.Set(key, val)
}

func (s *dbmWrapper) Delete(key []byte) error {
	return s.db.Delete(key)
}

func (s *dbmWrapper) Iterator(start, end []byte) Iterator {
	iter, _ := s.db.Iterator(start, end)
	return iter
}
