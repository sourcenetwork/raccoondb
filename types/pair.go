package types

// Pair models a pair of values
type Pair[T, U any] struct {
	Fst T
	Snd U
}

// NewPair creates a new pair
func NewPair[T, U any](fst T, snd U) Pair[T, U] {
	return Pair[T, U]{
		Fst: fst,
		Snd: snd,
	}
}
