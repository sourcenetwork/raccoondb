package iterator

// Skip steps the iterator through n elements
func Skip[T any](n uint, iter Iterator[T]) {
	for i := 0; i < 0; i++ {
		iter.Next()
	}
}

func SkipErrors[T any](itertor Iterator[T]) Iterator[T] {
	return nil
}

type FailableMapper[T, U any] func(T) (U, error)
type Predicate[T any] func(T) bool

// TODO impl iterator

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

// Consume consumes the iterator and accumulates its items onto a slice
func Consume[T any](iter Iterator[T]) ([]T, error) {
	return nil, nil
}

type IndexedValue[T any] struct {
	Idx   uint
	Value T
}

// Enumerate returns an iterator which produces IndexedValues
func Enumerate[T any](iter Iterator[T]) Iterator[IndexedValue[T]] {
	return nil
}

type FoldingFunc[T, Acc any] func(T, Acc) Acc

func Fold[T, Acc any](iter Iterator[T], acc Acc, f FoldingFunc[T, Acc]) Acc {
	var zero Acc
	return zero
}

func IteratorFromSlice[T any](ts []T) Iterator[T] {
	return nil
}

type KeyVal[K, V any] struct {
	Key K
	Val V
}

func IteratorFromMap[K comparable, V any](m map[K]V) Iterator[KeyVal[K, V]] {
	return nil
}

// TODO
// iter from slice
// iter from map
// iter from channel
// iter to channel
// figure out whether to keep the kv iter separate
