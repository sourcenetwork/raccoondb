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

// EncodeInt64 converts a int64 into a ones complement representation with a sign bit where 0 is negative and 1 is positive.
// Stores the result in a byte slice using big endian ordering.
// The resulting byte representation is comparable wrt the original values
func EncodeInt64(i int64) []byte {
	var j uint64
	if i > 0 {
		j = uint64(i)
		j = j | (1 << 63) // sets the MSB of j as 1 for positive
	} else {
		//makes i positive, flips its bits and sets the MSB as 0 in order to make negative number smaller than positive ones
		j = ^uint64(-i) & ^uint64(1<<63)
	}
	return EncodeUInt(j)
}

// DecodeInt64 converts a byte slice of a signed prefixed int64 into its original value.
func DecodeInt64(bytes []byte) int64 {
	var i int64

	ui := DecodeUInt(bytes)
	if ui>>63 == 1 { // pos number since MSB is 1
		ui = ui & ^uint64(1<<63) // sets the MSB to 0
		i = int64(ui)
	} else {
		ui = ^(ui | uint64(1<<63))
		i = -int64(ui)
	}
	return i
}
