package raccoon

import (
	"testing"
        "github.com/stretchr/testify/assert"
        "github.com/cosmos/cosmos-sdk/store/mem"
)

// Simple implementation of Edge interface
type record struct {
    Data string
    Source uint8
    Dest uint8
}

func (r record) GetSource() uint8 {
    return r.Source
}

func (r record) GetDest() uint8 {
    return r.Dest
}

func newRecord(data string, source, dest uint8) record {
    return record {
        Data: data,
        Source: source,
        Dest: dest,
    }
}

type recordKeyer struct {}

func (k *recordKeyer) Key(node uint8) []byte {
    return []byte{node}
}

func (k *recordKeyer) MinKey() []byte {
    return []byte{0}
}

func (k *recordKeyer) MaxKey() []byte {
    return []byte{255}
}


var _ Edge[uint8] = (*record)(nil)
var _ NodeKeyer[uint8] = (*recordKeyer)(nil)


func TestEdgeStoreSetThenGet(t *testing.T) {
    kv := mem.NewStore()
    var keyer NodeKeyer[uint8] = &recordKeyer{}
    marshaler := &Json[record]{}
    var store EdgeStore[record, uint8] = NewEdgeStore[record, uint8](kv, nil, keyer, marshaler)

    record := newRecord("data", 1, 2)
    err := store.Set(record)

    assert.Nil(t, err)

    opt, err := store.GetEdg(1, 2)
    assert.Nil(t, err)
    assert.False(t, opt.IsEmpty())
    assert.Equal(t, opt.Value(), record)
}

func TestEdgeStoreCannotGetADeleteValue(t *testing.T) {
    kv := mem.NewStore()
    var keyer NodeKeyer[uint8] = &recordKeyer{}
    marshaler := &Json[record]{}
    var store EdgeStore[record, uint8] = NewEdgeStore[record, uint8](kv, nil, keyer, marshaler)

    record := newRecord("data", 1, 2)
    err := store.Set(record)
    assert.Nil(t, err)

    err = store.Delete(record)
    assert.Nil(t, err)

    opt, err := store.GetEdg(1, 2)
    assert.Nil(t, err)
    assert.True(t, opt.IsEmpty())
}

func TestEdgeStoreGetSucessorsReturnDirectSucessors(t *testing.T) {
    kv := mem.NewStore()
    var keyer NodeKeyer[uint8] = &recordKeyer{}
    marshaler := &Json[record]{}
    var store EdgeStore[record, uint8] = NewEdgeStore[record, uint8](kv, nil, keyer, marshaler)

    rec1 := newRecord("data", 1, 2)
    err := store.Set(rec1)
    assert.Nil(t, err)

    rec2 := newRecord("data", 1, 3)
    err = store.Set(rec2)
    assert.Nil(t, err)

    sucessors, err := store.GetSucessors(1)
    assert.Nil(t, err)
    assert.Equal(t, 2, len(sucessors))
    assert.Contains(t, sucessors, rec1)
    assert.Contains(t, sucessors, rec2)
}

func TestEdgeStoreGetSucessorsReturnDirectAncestors(t *testing.T) {
    kv := mem.NewStore()
    var keyer NodeKeyer[uint8] = &recordKeyer{}
    marshaler := &Json[record]{}
    var store EdgeStore[record, uint8] = NewEdgeStore[record, uint8](kv, nil, keyer, marshaler)

    rec1 := newRecord("data", 1, 2)
    err := store.Set(rec1)
    assert.Nil(t, err)

    rec2 := newRecord("data", 3, 2)
    err = store.Set(rec2)
    assert.Nil(t, err)

    rec3 := newRecord("data", 9, 5)
    err = store.Set(rec3)
    assert.Nil(t, err)

    sucessors, err := store.GetAncestors(2)
    assert.Nil(t, err)
    assert.Equal(t, 2, len(sucessors))
    assert.Contains(t, sucessors, rec1)
    assert.Contains(t, sucessors, rec2)
}

func TestSecondaryIdxStoreSetThenGet(t *testing.T) {
    kv := mem.NewStore()
    var keyer NodeKeyer[uint8] = &recordKeyer{}
    marshaler := &Json[record]{}
    mapper := func(r record) []byte { return []byte(r.Data) }
    var store SecondaryIndexStore[record, uint8] = NewSecondaryIdxStore[record, uint8](kv, nil, "index", keyer, mapper, marshaler)

    rec1 := newRecord("data", 1, 2)
    err := store.Set(rec1)
    assert.Nil(t, err)

    rec2 := newRecord("data", 2, 3)
    err = store.Set(rec2)
    assert.Nil(t, err)

    records, err := store.GetByIdx([]byte("data"))
    assert.Nil(t, err)
    assert.Equal(t, 2, len(records))
    assert.Contains(t, records, rec1)
    assert.Contains(t, records, rec2)
}

func TestSecondaryIdxStoreSetThenGetFromUnknwonIdxReturnsEmptySlice(t *testing.T) {
    kv := mem.NewStore()
    var keyer NodeKeyer[uint8] = &recordKeyer{}
    marshaler := &Json[record]{}
    mapper := func(r record) []byte { return []byte(r.Data) }
    var store SecondaryIndexStore[record, uint8] = NewSecondaryIdxStore[record, uint8](kv, nil, "index", keyer, mapper, marshaler)

    rec1 := newRecord("data", 1, 2)
    err := store.Set(rec1)
    assert.Nil(t, err)

    rec2 := newRecord("data", 2, 3)
    err = store.Set(rec2)
    assert.Nil(t, err)

    println("BAD")
    records, err := store.GetByIdx([]byte("404"))
    assert.Nil(t, err)
    assert.Equal(t, 0, len(records))
}

func TestSecondaryIdxStoreDelete(t *testing.T) {
    kv := mem.NewStore()
    var keyer NodeKeyer[uint8] = &recordKeyer{}
    marshaler := &Json[record]{}
    mapper := func(r record) []byte { return []byte(r.Data) }
    var store SecondaryIndexStore[record, uint8] = NewSecondaryIdxStore[record, uint8](kv, nil, "index", keyer, mapper, marshaler)

    rec1 := newRecord("data", 1, 2)
    store.Set(rec1)
    rec2 := newRecord("data", 2, 3)
    store.Set(rec2)

    err := store.Delete(rec1)
    assert.Nil(t, err)

    records, err := store.GetByIdx([]byte("data"))
    assert.Nil(t, err)
    assert.Equal(t, 1, len(records))
    assert.NotContains(t, records, rec1)
    assert.Contains(t, records, rec2)
}
