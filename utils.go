package raccoon

import (
    "github.com/cosmos/cosmos-sdk/store/types"
)

func ProduceTrue[T any](t T) bool {
    return true
}

// Accumulates every item in iterator into a slice
func consumeIter[T any](iter types.Iterator, marshaler Marshaler[T]) ([]T, error) {
    return filterConsumeIter[T](iter, marshaler, ProduceTrue[T])
}

// Accumlates iterator items that match the predicate function into a slice.
func filterConsumeIter[T any](iter types.Iterator, marshaler Marshaler[T], predicate Predicate[T]) ([]T, error) {
    defer iter.Close()

    var records []T
    for ; iter.Valid(); iter.Next() {
        bytes := iter.Value()
        if err := iter.Error(); err != nil {
            return nil, err
        }

        record, err := marshaler.Unmarshal(bytes)
        if err != nil {
            return nil, err
        }

        if predicate(record) {
            records = append(records, record)
        }
    }
    return records, nil
}
