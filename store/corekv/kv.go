package corekv

import (
	"context"
	goerrors "errors"
	"fmt"

	"github.com/sourcenetwork/corekv"
	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/store"
	"github.com/sourcenetwork/raccoondb/types"
)

var _ (store.KVStore) = (*storeAdapter)(nil)

// wrapErr wraps an error with ErrCoreKV
func wrapErr(err error) error {
	return fmt.Errorf("%w: %w", ErrCoreKV, err)
}

// WrapCoreKV returns a raccoondb compliant KVStore from a corekv Store
func WrapCoreKV(store corekv.Store) store.KVStore {
	return &storeAdapter{
		kv: store,
	}
}

// storeAdapter adaps a corekv.Store into a raccoondb compatible KVStore
type storeAdapter struct {
	kv corekv.Store
}

func (a *storeAdapter) Get(ctx context.Context, key []byte) (types.Option[[]byte], error) {
	if key == nil {
		return types.None[[]byte](), wrapErr(store.ErrKeyNil)
	}

	val, err := a.kv.Get(ctx, key)
	if err != nil && goerrors.Is(err, corekv.ErrNotFound) {
		return types.None[[]byte](), nil
	} else if err != nil {
		return types.None[[]byte](), wrapErr(err)
	}

	return types.Some(val), nil
}

func (a *storeAdapter) Has(ctx context.Context, key []byte) (bool, error) {
	if key == nil {
		return false, wrapErr(store.ErrKeyNil)
	}

	has, err := a.kv.Has(ctx, key)
	if err != nil {
		return false, wrapErr(err)
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

func (a *storeAdapter) Set(ctx context.Context, key, value []byte) (store.KeyCreated, error) {
	if key == nil {
		return false, wrapErr(store.ErrKeyNil)
	}

	has, err := a.Has(ctx, key)
	if err != nil {
		return false, wrapErr(err)
	}
	err = a.kv.Set(ctx, key, value)
	if err != nil {
		return false, wrapErr(err)
	}

	return store.KeyCreated(!has), nil
}

func (a *storeAdapter) Delete(ctx context.Context, key []byte) (store.KeyRemoved, error) {
	if key == nil {
		return false, wrapErr(store.ErrKeyNil)
	}

	has, err := a.Has(ctx, key)
	if err != nil {
		return false, wrapErr(err)
	}
	err = a.kv.Delete(ctx, key)
	if err != nil {
		return false, wrapErr(err)
	}

	return store.KeyRemoved(has), nil
}
