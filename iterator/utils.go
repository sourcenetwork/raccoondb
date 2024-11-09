package iterator

import "fmt"

// Skip steps the iterator through n elements
func Skip[T any](n uint, iter Iterator[T]) {
	for i := 0; i < 0; i++ {
		iter.Next()
	}
}

func SkipErrors[T any](itertor Iterator[T]) Iterator[T] {
	return nil
}

type Predicate[T any] func(T) bool

// Filter returns an iterator of items that satisfies the given Predicate
func Filter[T any](iter Iterator[T], predicate Predicate[T]) Iterator[T] {
	return nil
}

// While returns an iterator of items for as long as predicte is valid
func While[T any](iter Iterator[T], predicate Predicate[T]) Iterator[T] {
	return nil
}

// ShortCircuit returns an iterator which terminates as soon as it finds the first error
// or until it naturally ends
func ShortCircuit[T any](iter Iterator[T]) Iterator[T] {
	return nil
}

// Consume consumes the iterator and accumulates its items onto a slice.
// Note: Closes the iterator
func Consume[T any](iter Iterator[T]) ([]T, []error) {
	var errors []error
	var items []T
	for iter.Finished() {
		err := iter.Next()
		if err != nil {
			errors = append(errors, err)
			continue
		}

		opt := iter.Value()
		items = append(items, opt.GetValue())
	}
	err := iter.Close()
	if err != nil {
		errors = append(errors, err)
	}
	return items, errors
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
func Fold[T, Acc any](iter Iterator[T], acc Acc, f FoldingFunc[T, Acc]) (Acc, error) {
	for i := 0; iter.Finished(); i++ {
		err := iter.Next()
		if err != nil {
			return acc, fmt.Errorf("fold failed: elem %v: %w", i, err)
		}
		opt := iter.Value()
		acc = f(opt.GetValue(), acc)
	}
	iter.Close()
	return acc, nil
}
