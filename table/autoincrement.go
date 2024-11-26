package table

import (
	"context"

	"github.com/sourcenetwork/raccoondb/primitives"
)

const counterKey string = "id"

type SetIDHook[T any] func(obj *T, id uint64) error

type AutoIncrementer[T any] struct {
	hook    SetIDHook[T]
	counter primitives.CounterStore
}

func (i *AutoIncrementer[T]) SetId(ctx context.Context, obj *T) (uint64, error) {
	id, err := i.counter.GetNext(ctx, []byte(counterKey))
	if err != nil {
		return 0, err
	}
	err = i.hook(obj, id)
	if err != nil {
		return 0, err
	}

	_, err = i.counter.Increment(ctx, []byte(counterKey))
	if err != nil {
		return 0, err
	}

	return id, nil
}
