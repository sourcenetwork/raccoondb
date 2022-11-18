package raccoon

import (
    "fmt"
)


// Separator used to build keys
const Separator = '/'


// Key represents a hiearachical key built from fragments.
// Key fragments are ordered and the final
// key is built joining fragments with a a constant separator.
type Key struct {
    fragments [][]byte
}

func (k Key) Fragments() [][]byte {
    return k.fragments
}

func (k Key) Append(fragments ...[]byte) Key {
    key := Key {
        fragments: make([][]byte, 0, len(k.fragments) + len(fragments)),
    }
    key.fragments = append(key.fragments, k.fragments...)
    key.fragments = append(key.fragments, fragments...)
    return key
}

func (k Key) Join(other Key) Key {
    key := Key{}
    key = key.Append(k.fragments...)
    key = key.Append(other.fragments...)
    return key
}

func (k Key) ToBytes() []byte {
    size := len(k.fragments)
    fmt.Printf("")
    for _, fragment := range k.fragments {
        size += len(fragment)
    }

    key := make([]byte, 0, size)

    for _, fragment := range k.fragments {
        key = append(key, byte(Separator))
        key = append(key, fragment...)
    }
    return key
}
