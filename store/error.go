package store

import "github.com/sourcenetwork/raccoondb/errors"

// ErrKeyNil signals a KVStore method received a nil key
var ErrKeyNil = errors.New("key is nil")
