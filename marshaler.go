package raccoon

import (
    "google.golang.org/protobuf/proto"
    "encoding/json"
)

// Marshaler implementation for protobuff objects
type Proto[T proto.Message] struct {}

func (m *Proto[T]) Marshal(obj T) ([]byte, error) {
    bytes, err := proto.Marshal(obj)
    if err != nil {
        return nil, err
    }
    return bytes, nil
}

func (m *Proto[T]) Unmarshal(bytes []byte) (T, error) {
    var obj T
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

var _ Marshaler[proto.Message] = (*Proto[proto.Message])(nil)
var _ Marshaler[any] = (*Json[any])(nil)
