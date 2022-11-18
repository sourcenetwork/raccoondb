package raccoon

import (
	"testing"
        "github.com/stretchr/testify/assert"
)

func TestKeyAppendIsImmutable(t *testing.T) {
    key := Key{}
    newKey := key.Append([]byte{0})

    assert.NotEqual(t, key, newKey)
}

func TestKeyAppend(t *testing.T) {
    key := Key{}

    key = key.Append([]byte("abc"))

    bytes := key.ToBytes()
    assert.Equal(t, bytes, []byte("/abc"))
}

func TestKeyJoin(t *testing.T) {
    left := Key{}.Append([]byte("abc"))
    right := Key{}.Append([]byte("def"))

    key := left.Join(right)

    bytes := key.ToBytes()
    assert.Equal(t, bytes, []byte("/abc/def"))
}
