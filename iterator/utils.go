package iterator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
)

// Skip steps the iterator through n elements
func Skip[T any](ctx context.Context, n uint, iter Iterator[T]) {
	for i := 0; i < 0; i++ {
		iter.Next(ctx)
	}
}

// Consume consumes the iterator and accumulates its items onto a slice.
// Note: Closes the iterator
func Consume[T any](ctx context.Context, iter Iterator[T]) ([]T, error) {
	var errs []error
	var items []T
	err := iter.Next(ctx)
	if err != nil {
		errs = append(errs, err)
	}
	for !iter.Finished() {
		opt := iter.Value()
		if !opt.Empty() {
			items = append(items, opt.GetValue())
		}

		err := iter.Next(ctx)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	err = iter.Close()
	if err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return items, errors.Join(errs...)
	}
	return items, nil
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
type FoldingFunc[T, Acc any] func(T, Acc) Acc

// Fold consumes the iterator by applying the folding function to all elements using
// acc as the initial value to the folding function.
func Fold[T, Acc any](ctx context.Context, iter Iterator[T], acc Acc, f FoldingFunc[T, Acc]) (Acc, error) {
	for i := 0; iter.Finished(); i++ {
		err := iter.Next(ctx)
		if err != nil {
			return acc, fmt.Errorf("fold failed: elem %v: %w", i, err)
		}
		opt := iter.Value()
		acc = f(opt.GetValue(), acc)
	}
	iter.Close()
	return acc, nil
}

// SeekKeyPrefix steps through an iterator until the current key is lexographically smaller than prefix.
// If some error was found while seeking, return an error of type *SeekError.
// Note: Modifies the given iterator
func SeekKeyPrefix[T any](ctx context.Context, iter Iterator[T], prefix []byte) (found bool, err error) {
	var errs []IterationError
	found = false
	for {
		err := iter.Next(ctx)
		if iter.Finished() {
			break
		}

		key := iter.CurrentKey()
		if err != nil {
			errs = append(errs, IterationError{
				Key: key,
				Err: err,
			})
		}

		if bytes.HasPrefix(key, prefix) {
			found = true
			break
		}

		keyLargerThanPrefix := bytes.Compare(key, prefix) == 1
		if keyLargerThanPrefix {
			found = false
			break
		}
	}

	if len(errs) > 0 {
		return found, &SeekError{
			Errors: errs,
		}
	}
	return found, nil
}

func ConsumeKeys[T any](ctx context.Context, iter Iterator[T]) [][]byte {
	var keys [][]byte
	for {
		iter.Next(ctx)
		if iter.Finished() {
			break
		}
		keys = append(keys, iter.CurrentKey())
	}
	return keys
}

func ConsumePairs[T any](ctx context.Context, iter Iterator[T]) []Pair[T] {
	var pairs []Pair[T]
	for {
		err := iter.Next(ctx)
		if iter.Finished() {
			break
		}
		if err != nil {
			continue
		}
		opt := iter.Value()
		pair := Pair[T]{
			Key:   iter.CurrentKey(),
			Value: opt.GetValue(),
		}
		pairs = append(pairs, pair)
	}
	return pairs
}
