package iterator

import (
	"bytes"
	"context"

	"github.com/sourcenetwork/raccoondb/marshal"
	"github.com/sourcenetwork/raccoondb/types"
	"github.com/sourcenetwork/raccoondb/utils"
)

var _ Iterator[any] = (*PairsIter[any])(nil)

// IterFromPairs returns an iterator which sorts paris then yields it sequentially
func IterFromPairs[T any](pairs []Pair[T]) Iterator[T] {
	sortable := utils.FromComparator(pairs, func(left, right Pair[T]) bool {
		return bytes.Compare(left.Key, right.Key) == -1
	})
	sortable.SortInPlace()

	return &PairsIter[T]{
		pairs: pairs,
		idx:   ^uint64(0),
		done:  false,
	}
}

// IterFromStringKeyMap returns an Iterator which steps through the pairs in the items map.
// The map keys are converted to bytes and ordered lexographically accodring to the bytes package.
func IterFromStringKeyMap[T any](items map[string]T) Iterator[T] {
	pairs := make([]Pair[T], 0, len(items))
	for key, value := range items {
		pair := Pair[T]{
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
	pairs := make([]Pair[T], 0, len(items))
	for i, item := range items {
		pair := Pair[T]{
			Key:   marshal.EncodeUInt(uint64(i)),
			Value: item,
		}
		pairs = append(pairs, pair)
	}
	return IterFromPairs(pairs)
}

func NewPair[T any](key []byte, val T) Pair[T] {
	return Pair[T]{
		Key:   key,
		Value: val,
	}
}

// Pair models a key value pair
type Pair[T any] struct {
	Key   []byte
	Value T
}

// PairIter represents an iterator which steps through a list of pairs
type PairsIter[T any] struct {
	pairs []Pair[T]
	idx   uint64
	done  bool
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

func (a *PairsIter[T]) Value() types.Option[T] {
	if a.done {
		return types.None[T]()
	}
	return types.Some(a.pairs[a.idx].Value)
}

func (a *PairsIter[T]) Finished() bool {
	return a.done
}

func (a *PairsIter[T]) Close() error {
	return nil
}

func (a *PairsIter[T]) GetParams() IteratorOpt {
	return IteratorOpt{
		start:   nil,
		end:     nil,
		prefix:  nil,
		reverse: false,
	}
}

func (a *PairsIter[T]) CurrentKey() []byte {
	if a.done {
		return nil
	}
	return a.pairs[a.idx].Key
}
