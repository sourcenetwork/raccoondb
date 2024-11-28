package table

import (
	"context"

	"github.com/sourcenetwork/raccoondb/marshal"
	"github.com/sourcenetwork/raccoondb/primitives"
	"github.com/sourcenetwork/raccoondb/store"
	"github.com/sourcenetwork/raccoondb/types"
)

const counterKey string = "id"

type IDSetter[T any] func(obj *T, id uint64)
type IDGetter[T any] func(obj *T) uint64

type AutoincrementTable[T any] struct {
	*Table[T]
	counter primitives.CounterStore
	Setter  IDSetter[T]
	Getter  IDGetter[T]
}

// Insert adds obj to the table.
// Fetches the next free ID from the table counter and sets it in obj
func (t *AutoincrementTable[T]) Insert(ctx context.Context, obj *T) error {
	id, err := t.counter.GetNext(ctx, []byte(counterKey))
	if err != nil {
		return err
	}
	t.Setter(obj, id)

	_, err = t.Table.Set(ctx, marshal.EncodeUInt(id), obj)
	if err != nil {
		return err
	}

	_, err = t.counter.Increment(ctx, []byte(counterKey))
	if err != nil {
		return err
	}

	return nil
}

func (t *AutoincrementTable[T]) GetByID(ctx context.Context, id uint64) (types.Option[T], error) {
	opt, err := t.Table.Get(ctx, []byte(counterKey))
	if err != nil {
		return types.None[T](), err
	}
	return opt, nil
}

func (t *AutoincrementTable[T]) DeleteByID(ctx context.Context, id uint64) (store.KeyRemoved, error) {
	removed, err := t.Table.Delete(ctx, marshal.EncodeUInt(id))
	if err != nil {
		return false, err
	}
	return removed, nil
}
