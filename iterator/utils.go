package iterator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
)

// Consume consumes the iterator and accumulates its items onto a slice.
// Note: Closes the iterator
func Consume[T any](ctx context.Context, iter Iterator[T]) ([]T, error) {
	var errs []error
	var items []T
	for !iter.Finished() {
		opt, err := iter.Value()
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if !opt.Empty() {
			items = append(items, opt.GetValue())
		}

		err = iter.Next(ctx)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	err := iter.Close()
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
		opt, err := iter.Value()
		if err != nil {
			return acc, fmt.Errorf("fold failed: elem %v: %w", i, err)
		}
		acc = f(opt.GetValue(), acc)

		err = iter.Next(ctx)
		if err != nil {
			return acc, fmt.Errorf("fold failed: elem %v: %w", i, err)
		}
	}
	iter.Close()
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

func ConsumePairs[T any](ctx context.Context, iter Iterator[T]) ([]Pair[T], error) {
	var pairs []Pair[T]
	var errs []IterItemError
	for {
		if iter.Finished() {
			break
		}

		opt, err := iter.Value()
		if err != nil {
			errs = append(errs, IterItemError{
				Key: iter.CurrentKey(),
				Err: err,
			})
		}
		if opt.Empty() {
			panic("opt empty")  TODO Fix this
		}

		pair := Pair[T]{
			Key:   iter.CurrentKey(),
			Value: opt.GetValue(),
		}
		pairs = append(pairs, pair)

		err = iter.Next(ctx)
		if err != nil {
			errs = append(errs, IterItemError{
				Key: iter.CurrentKey(),
				Err: err,
			})
		}
	}

	if len(errs) > 0 {
		return pairs, &IterationError{Errors: errs}
	}

	return pairs, nil
}
