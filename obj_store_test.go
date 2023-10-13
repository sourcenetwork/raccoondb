package raccoon

import (
    "testing"

        "github.com/stretchr/testify/require"
)

type dataProtoIder  struct {}
func (i *dataProtoIder) Id(o *Data) []byte {
    return []byte(o.Data)
}

func TestListObjIdsReturnsAllIdsInStore(t *testing.T) {
    kv := NewMemKV()
    var marshaler Marshaler[*Data] = &Json[*Data]{}
    store := NewObjStore[*Data](kv, marshaler, &dataProtoIder{})

    d1 := Data{Data: "d1"}
    err := store.SetObject(&d1)
    require.Nil(t, err)

    d2 := Data{Data: "d2"}
    err = store.SetObject(&d2)
    require.Nil(t, err)


    got, err := store.ListIds()
    want := [][]byte{
        []byte("d1"),
        []byte("d2"),
    }

    require.Nil(t, err)
    require.Equal(t, want, got)
}
