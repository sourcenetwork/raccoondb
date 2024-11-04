package types

// Marshaler marshalls and unmarshalls objects into a byte array,
// which gets persisted in the underlying store
type Marshaler[T any] interface {
	Marshal(*T) ([]byte, error)
	Unmarshal([]byte) (T, error)
}
