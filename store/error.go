package store

import "github.com/sourcenetwork/raccoondb/v2/errors"

// ErrKeyNil signals a KVStore method received a nil key
var ErrKeyNil = errors.New("key is nil")
