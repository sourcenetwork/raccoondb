package types

// Option represents a container for a value, which may or may not be empty
type Option[T any] struct {
	value T
	empty bool
}

// Empty returns true if the Option contains no value
func (o *Option[T]) Empty() bool {
	return o.empty
}

// GetValue returns value contained in the option.
// Panics if the Option is empty
func (o *Option[T]) GetValue() T {
	if o.empty {
		panic("option is empty")
	}

	return o.value
}

// None returns an empty Option
func None[T any]() Option[T] {
	return Option[T]{
		empty: true,
	}
}

// Some returns a new Option containing val
func Some[T any](val T) Option[T] {
	return Option[T]{
		value: val,
		empty: false,
	}
}
