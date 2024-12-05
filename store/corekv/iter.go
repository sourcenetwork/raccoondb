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
	iter        corekv.Iterator
	params      store.IterationParam
	initialized bool
}

func (i *iterAdapter) Next(_ context.Context) error {
	if !i.initialized {
		i.initialized = true
		return nil
	}

	i.iter.Next()
	return nil
}

func (i *iterAdapter) Value() types.Option[[]byte] {
	if i.Finished() || !i.initialized {
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

func (i *iterAdapter) GetParams() store.IterationParam {
	return i.params
}

func (i *iterAdapter) CurrentKey() []byte {
	if i.Finished() || !i.initialized {
		return nil
	}
	return i.iter.Key()
}
