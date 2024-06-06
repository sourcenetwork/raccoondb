package raccoon

import "fmt"

func NewPersistentKV(path, file string) (KVStore, error) {
	return nil, fmt.Errorf("Persistent Store not supported in WASM")
}
