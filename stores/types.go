package stores

import (
	"context"

	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/types"
)

type RecordRemoved bool
type RecordCreated bool

type ReadStore interface {
	Get(ctx context.Context, key []byte) (types.Option[[]byte], error)
	Has(ctx context.Context, key []byte) (bool, error)
	Iterate(ctx context.Context, opt iterator.IteratorOpt) (iterator.BytesIterator, error)
}

type KVStore interface {
	ReadStore
	Set(ctx context.Context, key, value []byte) (RecordCreated, error)
	Delete(ctx context.Context, key []byte) (RecordRemoved, error)
}
