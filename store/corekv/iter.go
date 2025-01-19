package corekv

import (
	"context"

	"github.com/sourcenetwork/corekv"
	"github.com/sourcenetwork/raccoondb/v2/store"
	"github.com/sourcenetwork/raccoondb/v2/types"
)

var _ (store.StoreIterator[[]byte]) = (*iterAdapter)(nil)

// iterAdapter adapts a corekv Iterator into a racoon iterator
type iterAdapter struct {
	iter     corekv.Iterator
	params   store.IterationParam
	finished bool
}

func (i *iterAdapter) Next(_ context.Context) error {
	i.iter.Next()
	if !i.iter.Valid() {
		i.finished = true
	}
	return nil
}

func (i *iterAdapter) Value() (types.Option[[]byte], error) {
	if i.finished {
		return types.None[[]byte](), nil
	}

	bytes, err := i.iter.Value()
	if err != nil {
		return types.None[[]byte](), err
	}

	if bytes == nil {
		return types.None[[]byte](), nil
	}
	return types.Some(bytes), nil
}

func (i *iterAdapter) Finished() bool {
	return i.finished
}

func (i *iterAdapter) Close() error {
	err := i.iter.Close(context.Background())
	if err != nil {
		return wrapErr(err)
	}
	return nil
}

func (i *iterAdapter) GetParams() store.IterationParam {
	return i.params
}

func (i *iterAdapter) CurrentKey() []byte {
	if i.finished {
		return nil
	}
	return i.iter.Key()
}
