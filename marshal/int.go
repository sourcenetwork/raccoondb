package marshal

import "encoding/binary"

// encode maps a uin64 to a big endian byte slice
func EncodeUInt(counter uint64) []byte {
	buff := make([]byte, 8)
	binary.BigEndian.PutUint64(buff, counter)
	return buff
}

// decode converts a big endian byte slice into a uint64
func DecodeUInt(counter []byte) uint64 {
	return binary.BigEndian.Uint64(counter)
}
