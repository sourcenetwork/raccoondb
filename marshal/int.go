package marshal

import (
	"encoding/binary"

	"github.com/sourcenetwork/raccoondb/errors"
)

var _ Marshaler[uint64] = UIntMarshaler{}

type UIntMarshaler struct{}

func (m UIntMarshaler) Marshal(i *uint64) ([]byte, error) {
	return EncodeUInt(*i), nil
}

func (m UIntMarshaler) Unmarshal(bytes []byte) (uint64, error) {
	if len(bytes) != 8 {
		return 0, errors.New("UIntMarshaler expected 8 byte slices")
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
