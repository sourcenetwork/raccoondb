package marshal

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/sourcenetwork/raccoondb/v2/errors"
)

var _ Marshaler[uint64] = UIntMarshaler{}
var _ Marshaler[string] = StringMarshaler{}
var _ Marshaler[[]byte] = BytesMarshaler{}
var _ Marshaler[any] = JsonMarshaler[any]{}

type UIntMarshaler struct{}

func (m UIntMarshaler) Marshal(i *uint64) ([]byte, error) {
	return EncodeUInt(*i), nil
}

func (m UIntMarshaler) Unmarshal(bytes []byte) (uint64, error) {
	if len(bytes) != 8 {
		return 0, errors.New("UIntMarshaler expects 8 byte slices")
	}

	return DecodeUInt(bytes), nil
}

// EncodeUInt maps a uin64 to a big endian byte slice
func EncodeUInt(i uint64) []byte {
	buff := make([]byte, 8)
	binary.BigEndian.PutUint64(buff, i)
	return buff
}

// DecodeUInt converts a big endian byte slice into a uint64
func DecodeUInt(bytes []byte) uint64 {
	return binary.BigEndian.Uint64(bytes)
}

type StringMarshaler struct{}

func (m StringMarshaler) Marshal(str *string) ([]byte, error) {
	return []byte(*str), nil
}

func (m StringMarshaler) Unmarshal(bz []byte) (string, error) {
	return string(bz), nil
}

func NewJSONMarshaler[T any](factory func() T) JsonMarshaler[T] {
	return JsonMarshaler[T]{
		factory: factory,
	}
}

type JsonMarshaler[T any] struct {
	factory func() T
}

func (m JsonMarshaler[T]) Marshal(obj *T) ([]byte, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("%w: json marshal: %w", ErrMarshaler, err)
	}
	return bytes, nil
}

func (m JsonMarshaler[T]) Unmarshal(bytes []byte) (T, error) {
	var zero T
	obj := m.factory()
	err := json.Unmarshal(bytes, &obj)
	if err != nil {
		return zero, fmt.Errorf("%w: json unmarshal: %w", ErrMarshaler, err)
	}
	return obj, nil
}

// BytesMarshaler implements Marshaler interface for byte slices
// Acts as an identity function
type BytesMarshaler struct{}

func (m BytesMarshaler) Marshal(bytes *[]byte) ([]byte, error) {
	return *bytes, nil
}

func (m BytesMarshaler) Unmarshal(bytes []byte) ([]byte, error) {
	return bytes, nil
}
