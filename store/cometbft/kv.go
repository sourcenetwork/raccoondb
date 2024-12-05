package cometbft

import (
	"context"

	cmdb "github.com/cometbft/cometbft-db"
	"github.com/sourcenetwork/raccoondb/v2/store"
	"github.com/sourcenetwork/raccoondb/v2/types"
)

var _ store.KVStore = (*kvWrapper)(nil)

// KVFromCometDB returns a Raccoon KVStore from a Comet DB
func KVFromCometDB(db cmdb.DB) store.KVStore {
	return &kvWrapper{
		db: db,
	}
}

// kvWrapper wraps a cometbft DB into a raccon KVStore
type kvWrapper struct {
	db cmdb.DB
}

func (k *kvWrapper) Get(ctx context.Context, key []byte) (types.Option[[]byte], error) {
	if key == nil {
		return types.None[[]byte](), wrapErr(store.ErrKeyNil)
	}

	bytes, err := k.db.Get(key)
	if err != nil {
		return types.None[[]byte](), wrapErr(err)
	}
	if bytes == nil {
		return types.None[[]byte](), nil
	}
	return types.Some(bytes), nil
}

func (k *kvWrapper) Has(ctx context.Context, key []byte) (bool, error) {
	if key == nil {
		return false, wrapErr(store.ErrKeyNil)
	}

	has, err := k.db.Has(key)
	if err != nil {
		return false, wrapErr(err)
	}
	return has, nil
}

func (k *kvWrapper) Iterate(ctx context.Context, opt store.IterationParam) (store.StoreIterator[[]byte], error) {
	var iter cmdb.Iterator
	var err error
	if opt.IsReverse() {
		iter, err = k.db.ReverseIterator(opt.GetLeftBound(), opt.GetRightBound())
	} else {
		iter, err = k.db.Iterator(opt.GetLeftBound(), opt.GetRightBound())
	}

	if err != nil {
		return nil, wrapErr(err)
	}
	wrapped := newWrappedIter(iter)
	return wrapped, nil
}

func (k *kvWrapper) Set(ctx context.Context, key, value []byte) (store.KeyCreated, error) {
	if key == nil {
		return false, wrapErr(store.ErrKeyNil)
	}

	has, err := k.db.Has(key)
	if err != nil {
		return false, wrapErr(err)
	}

	err = k.db.Set(key, value)
	if err != nil {
		return false, wrapErr(err)
	}
	return store.KeyCreated(!has), nil
}

func (k *kvWrapper) Delete(ctx context.Context, key []byte) (store.KeyRemoved, error) {
	if key == nil {
		return false, wrapErr(store.ErrKeyNil)
	}

	has, err := k.db.Has(key)
	if err != nil {
		return false, wrapErr(err)
	}

	err = k.db.Delete(key)
	if err != nil {
		return false, wrapErr(err)
	}
	return store.KeyRemoved(has), nil
}
