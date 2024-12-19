package table

import (
	"context"

	"github.com/sourcenetwork/raccoondb/v2/marshal"
	"github.com/sourcenetwork/raccoondb/v2/primitives"
	"github.com/sourcenetwork/raccoondb/v2/store"
	"github.com/sourcenetwork/raccoondb/v2/types"
)

func NewAutoIncrementer[T any](t *Table[T], getter IDGetter[T], setter IDSetter[T]) *Autoincrementer[T] {
	counterKv := primitives.NewPrefixedKV(t.baseStore, []byte(counterPrefix))
	counter := primitives.NewCounterStore(counterKv)
	return &Autoincrementer[T]{
		table:   t,
		counter: counter,
		setter:  setter,
		getter:  getter,
	}
}

const counterPrefix string = "counter/"
const counterKey string = "id"

// IDSetter models a hook which is used by AutoincrementTable
// to set the new ID of an object beign inserted
type IDSetter[T any] func(obj *T, id uint64)

// IDGetter models a hook which is used by AutoincrementTable
// to get the existing ID from the given object
type IDGetter[T any] func(obj *T) uint64

// AutoIncrementTable adds a top level counter to a Table,
// which is used to generate Identifiers for objects stored in Table.
//
// Identifiers are unsigned integers, which are marshaled using big endian encoding.
type Autoincrementer[T any] struct {
	table   *Table[T]
	counter primitives.CounterStore
	setter  IDSetter[T]
	getter  IDGetter[T]
}

// Insert adds obj to the table.
// Fetches the next free ID from the table counter and sets it in obj
func (t *Autoincrementer[T]) Insert(ctx context.Context, obj *T) error {
	id, err := t.counter.GetNext(ctx, []byte(counterKey))
	if err != nil {
		return err
	}
	t.setter(obj, id)

	_, err = t.table.Set(ctx, marshal.EncodeUInt(id), *obj)
	if err != nil {
		return err
	}

	_, err = t.counter.Increment(ctx, []byte(counterKey))
	if err != nil {
		return err
	}

	return nil
}

// Update modifies record in the table.
// Fetches the current ID from obj and updates the record with the recovered ID.
func (t *Autoincrementer[T]) Update(ctx context.Context, obj *T) error {
	id := t.getter(obj)
	_, err := t.table.Set(ctx, marshal.EncodeUInt(id), *obj)
	if err != nil {
		return err
	}
	return nil
}

// GetByID returns the object stored with the given integer id
func (t *Autoincrementer[T]) GetByID(ctx context.Context, id uint64) (types.Option[T], error) {
	opt, err := t.table.Get(ctx, marshal.EncodeUInt(id))
	if err != nil {
		return types.None[T](), err
	}
	return opt, nil
}

// GetByID returns the object stored with the given integer id
func (t *Autoincrementer[T]) GetByRecordID(ctx context.Context, record *T) (types.Option[T], error) {
	intId := t.getter(record)
	opt, err := t.GetByID(ctx, intId)
	if err != nil {
		return types.None[T](), err
	}
	return opt, nil
}

// DeleteByID removes the object stored with the given integer id
func (t *Autoincrementer[T]) DeleteByID(ctx context.Context, id uint64) (store.KeyRemoved, error) {
	removed, err := t.table.Delete(ctx, marshal.EncodeUInt(id))
	if err != nil {
		return false, err
	}
	return removed, nil
}
