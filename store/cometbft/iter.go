package cometbft

import (
	"context"

	cmdb "github.com/cometbft/cometbft-db"

	"github.com/sourcenetwork/raccoondb/v2/store"
	"github.com/sourcenetwork/raccoondb/v2/types"
)

var _ store.StoreIterator[[]byte] = (*iterWrapper)(nil)

func newWrappedIter(iter cmdb.Iterator) store.StoreIterator[[]byte] {
	return &iterWrapper{
		i:        iter,
		finished: false,
	}
}

type iterWrapper struct {
	i        cmdb.Iterator
	params   store.IterationParam
	finished bool
}

func (i *iterWrapper) Next(ctx context.Context) error {
	if i.finished {
		return nil
	}
	i.i.Next()

	if !i.i.Valid() {
		i.finished = true
	}

	err := i.i.Error()
	if err != nil {
		return wrapErr(err)
	}

	return nil
}

func (i *iterWrapper) Value() (types.Option[[]byte], error) {
	if i.finished {
		return types.None[[]byte](), nil
	}
	return types.Some(i.i.Value()), nil
}

func (i *iterWrapper) Finished() bool {
	return i.finished
}

func (i *iterWrapper) Close() error {
	err := i.i.Close()
	if err != nil {
		return wrapErr(err)
	}
	return nil
}

func (i *iterWrapper) GetParams() store.IterationParam {
	return i.params
}

func (i *iterWrapper) CurrentKey() []byte {
	if i.finished {
		return nil
	}
	return i.i.Key()
}
