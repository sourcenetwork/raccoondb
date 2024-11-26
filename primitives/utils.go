package primitives

// concatKey returns a slice which contains the concatination of prefix with key
func concatKey(prefix, key []byte) []byte {
	bytes := make([]byte, 0, len(prefix)+len(key))
	bytes = append(bytes, prefix...)
	bytes = append(bytes, key...)
	return bytes
}
