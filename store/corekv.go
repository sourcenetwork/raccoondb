package store

import (
	"context"
	"errors"

	"github.com/sourcenetwork/corekv"
	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/types"
)

var _ (iterator.Iterator[[]byte]) = (*iterAdapter)(nil)
var _ (KVStore) = (*storeAdapter)(nil)

// WrapCoreKV returns a raccoondb compliant KVStore from a corekv Store
func WrapCoreKV(store corekv.Store) KVStore {
	return &storeAdapter{
		kv: store,
	}
}

type storeAdapter struct {
	kv corekv.Store
}

func (a *storeAdapter) Get(ctx context.Context, key []byte) (types.Option[[]byte], error) {
	val, err := a.kv.Get(ctx, key)
	if err != nil && errors.Is(err, corekv.ErrNotFound) {
		return types.None[[]byte](), nil
	} else if err != nil {
		return types.None[[]byte](), err // TODO wrap
	}
	return types.Some(val), nil
}

func (a *storeAdapter) Has(ctx context.Context, key []byte) (bool, error) {
	has, err := a.kv.Has(ctx, key)
	if err != nil {
		return false, err //TODO WRAP
	}
	return has, nil
}

func (a *storeAdapter) Iterate(ctx context.Context, opt iterator.IteratorOpt) (iterator.BytesIterator, error) {
	// TODO fix opts
	iter := a.kv.Iterator(ctx, corekv.IterOptions{})
	return &iterAdapter{
		iter: iter,
	}, nil
}

func (a *storeAdapter) Set(ctx context.Context, key, value []byte) (RecordCreated, error) {
	has, err := a.Has(ctx, key)
	if err != nil {
		return false, err
	}
	err = a.kv.Set(ctx, key, value)
	if err != nil {
		return false, err
	}

	if has {
		return false, nil
	}

	return true, nil
}

func (a *storeAdapter) Delete(ctx context.Context, key []byte) (RecordRemoved, error) {
	has, err := a.Has(ctx, key)
	if err != nil {
		return false, err
	}
	err = a.kv.Delete(ctx, key)
	if err != nil {
		return false, err
	}

	if has {
		return false, nil
	}

	return true, nil
}

type iterAdapter struct {
	iter   corekv.Iterator
	params iterator.IteratorOpt
}

func (i *iterAdapter) Next() error {
	i.iter.Next()
	return nil
}

func (i *iterAdapter) Value() types.Option[[]byte] {
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
		return err // WRAP
	}
	return nil
}

func (i *iterAdapter) GetParams() iterator.IteratorOpt {
	return i.params
}

func (i *iterAdapter) CurrentKey() []byte {
	return i.iter.Key()
}
