package iterator

type KeyVal[K, V any] struct {
	Key K
	Val V
}

func IteratorFromMap[K comparable, V any](m map[K]V) Iterator[KeyVal[K, V]] {
	return nil
}
