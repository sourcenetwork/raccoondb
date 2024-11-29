package types

type Pair[T, U any] struct {
	Fst T
	Snd U
}

func NewPair[T, U any](fst T, snd U) Pair[T, U] {
	return Pair[T, U]{
		Fst: fst,
		Snd: snd,
	}
}
