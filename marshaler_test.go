package raccoon

import (
    "testing"

    "github.com/stretchr/testify/assert"

)

func TestProtoMarshaler(t *testing.T) {
    d := Data{
        Data: "abc",
    }

    m := ProtoMarshaler[*Data](func () *Data {return &Data{}})

    bytes, err := m.Marshal(&d)
    assert.Nil(t, err)
    assert.NotNil(t, bytes)

    got, err := m.Unmarshal(bytes)
    assert.Nil(t, err)
    assert.Equal(t, got.Data, "abc")
}
