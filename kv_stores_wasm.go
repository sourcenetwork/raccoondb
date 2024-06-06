package raccoon

import "fmt"

func NewPersistentKV(path, file string) (KVStore, CleanupFn, error) {
	return nil, nil, fmt.Errorf("Persistent Store not supported in WASM")
}
