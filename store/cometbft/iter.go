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
		i:           iter,
		initialized: false,
		finished:    false,
	}
}

type iterWrapper struct {
	i           cmdb.Iterator
	params      store.IteratorOpt
	initialized bool
	finished    bool
}

func (i *iterWrapper) Next(ctx context.Context) error {
	if !i.initialized {
		i.initialized = true
		// cometbft-db's iterator is created ready to use (yields first value right away)
		// as such it may have an error set during creation
		// if it fails to yield the first value, therefore we check for it
		if i.i.Error() != nil {
			return wrapErr(i.i.Error())
		}
		return nil
	}

	i.i.Next()
	err := i.i.Error()
	if !i.i.Valid() {
		i.finished = true
	}

	if err != nil {
		return wrapErr(err)
	}
	return nil
}

func (i *iterWrapper) Value() types.Option[[]byte] {
	if i.finished || !i.initialized {
		return types.None[[]byte]()
	}
	return types.Some(i.i.Value())
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

func (i *iterWrapper) GetParams() store.IteratorOpt {
	return i.params
}

func (i *iterWrapper) CurrentKey() []byte {
	if i.finished || !i.initialized {
		return nil
	}
	return i.i.Key()
}
