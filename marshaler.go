package raccoon

import (
    "encoding/json"

    "google.golang.org/protobuf/proto"
)

type ProtoConstraint[T any] interface {
    proto.Message
    *T
}

/*
// Marshaler implementation for protobuff objects
type ProtoMarshaler[T any, PT ProtoConstraint[T]] struct {}


func (m *ProtoMarshaler[T, PT]) Marshal(obj PT) ([]byte, error) {
    bytes, err := proto.Marshal(obj)
    if err != nil {
        return nil, err
    }
    return bytes, nil
}

func (m *ProtoMarshaler[T, PT]) Unmarshal(bytes []byte) (T, error) {
    var obj T
    p := PT(&obj)
    
    err := proto.Unmarshal(bytes, p)
    if err != nil {
        return obj, err
    }
    return obj, nil
}
*/

// Marshaler implementation for protobuff objects
type factoryProtoMarshaler[T proto.Message] struct {
    factory func() T
}

func ProtoMarshaler[T proto.Message](factory func() T) Marshaler[T] {
    return &factoryProtoMarshaler[T]{
        factory: factory,
    }
}

func (m *factoryProtoMarshaler[T]) Marshal(obj *T) ([]byte, error) {
    bytes, err := proto.Marshal(*obj)
    if err != nil {
        return nil, err
    }
    return bytes, nil
}

func (m *factoryProtoMarshaler[T]) Unmarshal(bytes []byte) (T, error) {
    obj := m.factory()
    err := proto.Unmarshal(bytes, obj)
    if err != nil {
        return obj, err
    }
    return obj, nil
}


// Marshaler to represent objects as Json
type Json[T any] struct {}

func (m *Json[T]) Marshal(obj *T) ([]byte, error) {
    bytes, err := json.Marshal(obj)
    if err != nil {
        return nil, err
    }
    return bytes, nil
}

func (m *Json[T]) Unmarshal(bytes []byte) (T, error) {
    var obj T
    err := json.Unmarshal(bytes, &obj)
    if err != nil {
        return obj, err
    }
    return obj, nil
}

//var _ Marshaler[Data, *Data] = (*ProtoMarshaler[Data, *Data])(nil)
var _ Marshaler[any] = (*Json[any])(nil)
