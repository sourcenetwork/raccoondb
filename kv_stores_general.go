//go:build !wasm

package raccoon

import (
	dbm "github.com/cosmos/cosmos-db"
)

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

func NewPersistentKV(path, file string) (KVStore, CleanupFn, error) {
	db, err := dbm.NewGoLevelDB(file, path, nil)
	if err != nil {
		return nil, nil, err
	}
	wrapper := dbmWrapper{
		db: db,
	}
	cleanup := func() error {
		return db.Close()
	}

	return &wrapper, cleanup, nil
}
