package raccoon

import (
    "encoding/json"

    "google.golang.org/protobuf/proto"
)

// Marshaler implementation for protobuff objects
type protoMarshaler[T proto.Message] struct {
    factory func() T
}

func ProtoMarshaler[T proto.Message](factory func() T) Marshaler[T]{
    return &protoMarshaler[T] {
        factory: factory,
    }
}

func (m *protoMarshaler[T]) Marshal(obj T) ([]byte, error) {
    bytes, err := proto.Marshal(obj)
    if err != nil {
        return nil, err
    }
    return bytes, nil
}

func (m *protoMarshaler[T]) Unmarshal(bytes []byte) (T, error) {
    obj := m.factory()
    
    err := proto.Unmarshal(bytes, obj)
    if err != nil {
        return obj, err
    }
    return obj, nil
}


// Marshaler to represent objects as Json
type Json[T any] struct {}

func (m *Json[T]) Marshal(obj T) ([]byte, error) {
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

var _ Marshaler[proto.Message] = (*protoMarshaler[proto.Message])(nil)
var _ Marshaler[any] = (*Json[any])(nil)
