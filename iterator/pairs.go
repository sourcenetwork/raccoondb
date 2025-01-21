package iterator

import (
	"bytes"
	"context"

	"github.com/sourcenetwork/raccoondb/v2/marshal"
	"github.com/sourcenetwork/raccoondb/v2/utils"
)

var _ Iterator[any] = (*PairsIter[any])(nil)

// IterFromPairs returns an iterator which sorts paris then yields it sequentially
func IterFromPairs[T any](pairs []KeyValue[T]) Iterator[T] {
	sortable := utils.FromComparator(pairs, func(left, right KeyValue[T]) bool {
		return bytes.Compare(left.Key, right.Key) == -1
	})
	sortable.SortInPlace()
	if len(pairs) == 0 {
		return NewEmptyIterator[T]()
	}

	return &PairsIter[T]{
		pairs: pairs,
		idx:   0,
		done:  false,
	}
}

// IterFromStringKeyMap returns an Iterator which steps through the pairs in the items map.
// The map keys are converted to bytes and ordered lexographically accodring to the bytes package.
func IterFromStringKeyMap[T any](items map[string]T) Iterator[T] {
	pairs := make([]KeyValue[T], 0, len(items))
	for key, value := range items {
		pair := KeyValue[T]{
			Key:   []byte(key),
			Value: value,
		}
		pairs = append(pairs, pair)
	}
	return IterFromPairs(pairs)
}

// IterFromslice returns an iterator which yields the values in items
// The iterator keys are the big endian encoding of the item's index in items
func IterFromSlice[T any](items []T) Iterator[T] {
	pairs := make([]KeyValue[T], 0, len(items))
	for i, item := range items {
		pair := KeyValue[T]{
			Key:   marshal.EncodeUInt(uint64(i)),
			Value: item,
		}
		pairs = append(pairs, pair)
	}
	return IterFromPairs(pairs)
}

func NewPair[T any](key []byte, val T) KeyValue[T] {
	return KeyValue[T]{
		Key:   key,
		Value: val,
	}
}

// PairIter represents an iterator which steps through a list of pairs
type PairsIter[T any] struct {
	pairs []KeyValue[T]
	idx   uint64
	done  bool
	zero  T
}

func (i *PairsIter[T]) Next(_ context.Context) error {
	if i.done {
		return nil
	}

	i.idx++
	if i.idx == uint64(len(i.pairs)) {
		i.done = true
		return nil
	}

	return nil
}

func (a *PairsIter[T]) Value() (T, error) {
	if a.done {
		return a.zero, nil
	}
	return a.pairs[a.idx].Value, nil
}

func (a *PairsIter[T]) Finished() bool {
	return a.done
}

func (a *PairsIter[T]) Close() error {
	return nil
}

func (a *PairsIter[T]) CurrentKey() []byte {
	if a.done {
		return nil
	}
	return a.pairs[a.idx].Key
}
