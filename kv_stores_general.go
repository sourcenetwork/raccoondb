//go:build !wasm

package raccoon

import (
	dbm "github.com/cosmos/cosmos-db"
)

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
