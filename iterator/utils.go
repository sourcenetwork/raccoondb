package iterator

import (
	"bytes"
	"context"
	"fmt"
)

// Consume consumes the iterator and accumulates its items into a slice.
// Note: Closes the iterator
func Consume[T any](ctx context.Context, iter Iterator[T]) ([]T, error) {
	var vals []T
	foldingFunc := func(kv KeyValue[T], acc []T) []T {
		return append(acc, kv.Value)
	}
	return Fold(ctx, iter, vals, foldingFunc)
}

// IndexedValue is a pair indicating the Value and the Idx that produced the element
type IndexedValue[T any] struct {
	Idx   uint
	Value T
}

// Enumerate returns an iterator which produces IndexedValues
func Enumerate[T any](iter Iterator[T]) Iterator[IndexedValue[T]] {
	var counter uint = 0
	ptr := &counter
	mapper := func(t T) IndexedValue[T] {
		val := IndexedValue[T]{
			Idx:   *ptr,
			Value: t,
		}
		*ptr = *ptr + 1
		return val
	}
	return Map(iter, mapper)
}

// FoldingFunc takes a value and an accumulator and returns an updated accumulator
type FoldingFunc[T, Acc any] func(KeyValue[T], Acc) Acc

// Fold consumes the iterator by applying the folding function to all elements using
// acc as the initial value to the folding function.
// Iterates until completion or until the first error
func Fold[T, Acc any](ctx context.Context, iter Iterator[T], acc Acc, f FoldingFunc[T, Acc]) (Acc, error) {
	defer iter.Close()
	for !iter.Finished() {
		val, err := iter.Value()
		if err != nil {
			return acc, fmt.Errorf("fold failed: key %v: %w", iter.CurrentKey(), err)
		}
		kv := KeyValue[T]{
			Key:   iter.CurrentKey(),
			Value: val,
		}
		acc = f(kv, acc)

		err = iter.Next(ctx)
		if err != nil {
			return acc, fmt.Errorf("fold failed: elem %v: %w", iter.CurrentKey(), err)
		}
	}
	return acc, nil
}

// SeekKeyPrefix steps through an iterator until the current key is lexographically smaller than prefix.
// Aggregates errors found in iteration and returns a *IterationError if any errors are found.
// Note: Modifies the given iterator
func SeekKeyPrefix[T any](ctx context.Context, iter Iterator[T], prefix []byte) (found bool, err error) {
	var errs []IterItemError
	found = false
	for {
		if iter.Finished() {
			break
		}

		key := iter.CurrentKey()
		if bytes.HasPrefix(key, prefix) {
			found = true
			break
		}

		keyLargerThanPrefix := bytes.Compare(key, prefix) == 1
		if keyLargerThanPrefix {
			found = false
			break
		}

		err := iter.Next(ctx)
		if err != nil {
			errs = append(errs, IterItemError{
				Key: key,
				Err: err,
			})
		}
	}

	if len(errs) > 0 {
		return found, &IterationError{Errors: errs}
	}
	return found, nil
}

// ConsumePairs steps through the iterator and returns all key-value pairs contained in it.
// Stops at the first error
func ConsumePairs[T any](ctx context.Context, iter Iterator[T]) ([]KeyValue[T], error) {
	var pairs []KeyValue[T]
	foldingFunc := func(kv KeyValue[T], acc []KeyValue[T]) []KeyValue[T] {
		return append(acc, kv)
	}
	return Fold(ctx, iter, pairs, foldingFunc)
}

// ConsumeKeys steps through the iterator and returns all keys contained in it.
// Stops at the first error
func ConsumeKeys[T any](ctx context.Context, iter Iterator[T]) ([][]byte, error) {
	var pairs [][]byte
	foldingFunc := func(kv KeyValue[T], acc [][]byte) [][]byte {
		return append(acc, kv.Key)
	}
	return Fold(ctx, iter, pairs, foldingFunc)
}
