package composite

/*

import (
	"context"

	"github.com/sourcenetwork/raccoondb/iterator"
	"github.com/sourcenetwork/raccoondb/marshal"
)

const autoIncCounterKey = "counter"

type PreSetCallback[O any] func(obj *O, id uint64)

func NewIDedObjectStore[Obj any](marshaler marshal.Marshaler[Obj], hook PreSetCallback[Obj]) IDedObjectStore[Obj] {
	return IDedObjectStore[Obj]{}
}

type IDedObjectStore[T any] struct {
	oidx    *IndexedObjectStore[T]
	counter *CounterStore
	hook    PreSetCallback[T]
}

func (s *IDedObjectStore[T]) Initialize(kv KVStore, indexes ...[]*Index[T, any]) error {
	// TODO create substores and so on
	return nil
}

func (s *IDedObjectStore[T]) Delete(ctx context.Context, id uint64) (RecordRemoved, error) {
	bz := encUint(id)
	return s.oidx.Delete(ctx, bz)
}

func (s *IDedObjectStore[T]) Set(ctx context.Context, obj *T) (RecordCreated, error) {
	// Fixme
	id, err := s.counter.GetNextAndIncrement(ctx, []byte(autoIncCounterKey))
	if err != nil {
		return false, err
	}
	s.hook(obj, id)

	return s.oidx.Set(ctx, encUint(id), obj)
}

func (s *IDedObjectStore[T]) UpdateIndexes(ctx context.Context) error {
	return s.UpdateIndexes(ctx)
}

func (s *IDedObjectStore[T]) Get(ctx context.Context, id uint64) (T, error) { var z T; return z, nil }

func (s *IDedObjectStore[T]) List(ctx context.Context) (iterator.Iterator[T], error) { return nil, nil }

func (s *IDedObjectStore[T]) Has(ctx context.Context, id uint64) (bool, error) { return false, nil }

*/
