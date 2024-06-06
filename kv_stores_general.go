//go:build !wasm

package raccoon

import (
	dbm "github.com/cosmos/cosmos-db"
)

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
