package corekv

import (
	"context"

	"github.com/sourcenetwork/corekv"
	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/types"
)

var _ (iterator.Iterator[[]byte]) = (*iterAdapter)(nil)

// iterAdapter adapts a corekv Iterator into a racoon iterator
type iterAdapter struct {
	iter   corekv.Iterator
	params iterator.IteratorOpt
}

func (i *iterAdapter) Next() error {
	i.iter.Next()
	return nil
}

func (i *iterAdapter) Value() types.Option[[]byte] {
	if i.Finished() {
		return types.None[[]byte]()
	}

	bytes := i.iter.Value()
	if bytes == nil {
		return types.None[[]byte]()
	}
	return types.Some(bytes)
}

func (i *iterAdapter) Finished() bool {
	return !i.iter.Valid()
}

func (i *iterAdapter) Close() error {
	err := i.iter.Close(context.Background())
	if err != nil {
		return wrapErr(err)
	}
	return nil
}

func (i *iterAdapter) GetParams() iterator.IteratorOpt {
	return i.params
}

func (i *iterAdapter) CurrentKey() []byte {
	if i.Finished() {
		return nil
	}
	return i.iter.Key()
}
