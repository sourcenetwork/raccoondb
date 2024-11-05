package stores

import (
	"context"

	"github.com/sourcenetwork/raccoondb/iterator"
)

// Index should group []byte *values* under []byte *buckets*
// Add typed index afterwards

const buckets = "buckets"
const idxPrefix = "idx"

type Indexable interface {
	ToBytes() []byte
}

func NewFieldIndexStore(kv KVStore) FieldIndexStore {
	return FieldIndexStore{
		baseKv:  kv,
		buckets: NewCountedKVStore(NewPrefixedKV(kv, []byte(buckets))),
		idx:     NewPrefixedKV(kv, []byte(idxPrefix)),
	}
}

type FieldIndexStore struct {
	baseKv  KVStore
	buckets CountedKVStore
	idx     KVStore
}

// FIXME this is not right, it doesn't overwrite / change a previous index
// it only appends
func (s *FieldIndexStore) IndexValue(ctx context.Context, bucket []byte, value []byte) error {
	_, err := s.buckets.Set(ctx, bucket, bucket)
	if err != nil {
		return err
	}

	bucketKey := getIdxKey(bucket, value)
	_, err = s.idx.Set(ctx, bucketKey, value)
	if err != nil {
		return err
	}

	return nil
}

func (s *FieldIndexStore) Has(ctx context.Context, bucket, value []byte) (bool, error) {
	idxKey := getIdxKey(bucket, value)
	return s.idx.Has(ctx, idxKey)
}

func (s *FieldIndexStore) GetBucketValues(ctx context.Context, bucket []byte) (iterator.BytesIterator, error) {
	prefixKv := NewPrefixedKV(s.idx, bucket)
	var opts iterator.IteratorOpt // FIXME add open iterator
	iter, err := prefixKv.Iterate(ctx, opts)
	if err != nil {
		return nil, err
	}
	return iter, nil
}

func (s *FieldIndexStore) Delete(ctx context.Context, value []byte) error {
	// FIXME - deleting an entry doesn't mean I should completely unindex a value
	// keep count of entries for bucket
	//err := s.vals.Delete(ctx, i.ToBytes())
	//if err != nil {
	//return err
	//}

	// FIXME need to change the delete to look for the current mapping
	// find the current key->value map and reverse it

	return nil
}

func (s *FieldIndexStore) IterIndexValues(ctx context.Context) (iterator.BytesIterator, error) {
	var opts iterator.IteratorOpt // FIXME add open iterator
	iter, err := s.buckets.Iterate(ctx, opts)
	if err != nil {
		return nil, err
	}
	return iter, nil
}

func (s *FieldIndexStore) GetCount(ctx context.Context) (uint64, error) {
	return s.buckets.GetCount(ctx)
}

func (s *FieldIndexStore) UpdateIndex(ctx context.Context, values iterator.BytesIterator) error {
	panic("todo")
}

func getIdxKey(bucket []byte, value []byte) []byte {
	l := len(bucket) + len(value) + 1
	key := make([]byte, 0, l)
	key = append(key, bucket...)
	key = append(key, byte('/'))
	key = append(key, value...)
	return key
}
