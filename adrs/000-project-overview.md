---
date: 2026-03-29
---

# RaccoonDB

RaccoonDB is a Go library which provides common primitives and abstractions for KV Stores, the goal is to simplify KV Store usage while building apps.

## Core Interfaces

RaccoonDB defines a layered interface hierarchy for store operations:

- **ReadStore**: The read-only foundation providing `Get`, `Has`, and `Iterate` operations
- **KVStore**: Extends ReadStore with mutation capabilities (`Set`, `Delete`)
- **CountedStore**: Provides count tracking via `GetCount`

This separation allows components to declare minimal required capabilities, enabling read-only access patterns where mutations aren't needed.

## Store Adapters

RaccoonDB is backend-agnostic. The library provides adapters for multiple underlying KV stores:

- **CoreKV**: Primary adapter wrapping `github.com/sourcenetwork/corekv`
- **CometBFT**: Adapter for CometBFT-DB (`github.com/cometbft/cometbft-db`)
- **Cosmos SDK**: Adapter for `cosmossdk.io/core/store.KVStore`

This allows applications to choose their storage backend based on their ecosystem (blockchain vs general purpose) while using the same RaccoonDB abstractions.

## Abstractions

### Marshalers
Raccoon has a built-in notion of Iterators because most of the time application are storing bytes which represent a serialized version of some domain object.
The Marshaler interface gives callers the flexibility to use their prefered marshalling scheme.

### Iterator
Iterators are a well known abstration in the KV Domain.
Raccoon takes Iterators a step futher and makes them somewhat composable, similar to Rust's Iter trait.

The iterator pkg provides utility functions to consume iterators as well as transform functions which wraps iterators, such as a Filter transformer which wraps an iterator and returns only items matching a given predicate.
Another good example is a Map transformer, which applies a mapping function to items in the iterator.

In general, the spirit of raccoon is to use a functional pattern around iterators to provide achieve a nice interface to deal with iterators.
The reason for that is twofold: this pattern allows callers to handle iterators lazily and it eliminates the needs for callbacks, which is a common design choice while dealing with iterators.

Example Iterator Usage:

```go
defer iter.Close()
for !iter.Finished() {
    _, err := iter.Value()
    if err != nil { return err }
    err := iter.Next();
    if err != nil { return err }
}
```

The Iterator abstraction differs a bit from the usual patterns in Go libraries, but it was designed to meet a few criteria:
1. The iterator is created "ready to use", no need to step before using.
   This choice was a consequence of the fact Raccoon extensively wraps iterators and this model solves some challenges with the transformers.
2. Value returns an error: also unorthodox, but because transformers can individually fail (such as mapping or marshaling a value can fail), we chose to separate a value fetching error from a 
   B-Tree / list / underlying store progression error
3. Related to 1, we wanted for the iterator to be ready to use, but not require a check or stepping through before the loop actually initializes, therefore the chice of `!iter.Finished()` as the loop
   condition requires no action on the users behalf when an iterator is possibly empty.

### Primitives

Primitives are the building blocks for application developers, meant to represent reusable and composable primitives over KVStores.

Key primitives include:

- **PrefixStore**: Wraps a KVStore with a key prefix, enabling logical namespacing without separate physical stores. This is used extensively throughout the library to organize data.
- **CounterStore**: Manages uint64 counters with increment/decrement operations, useful for tracking counts without full iteration.
- **CountedKVStore**: A KVStore that maintains an entry count, combining storage with automatic count tracking.
- **KeyObjectStore[T]**: Generic type-safe storage that handles marshaling/unmarshaling automatically, bridging raw bytes and typed domain objects.
- **FieldIndexStore**: Bucket-based indexing structure that groups items by field values, maintaining per-bucket counts for efficient queries.


### Table

The Table abstraction is a top level abstraction which composes several of the underlying primitives.
It's meant to model a data storage abstraction which has some of the niceties of a database table developers are used from the DB world.

Tables can have indices defined on them, synchronizes index removal and updating while creating objects, asserting there is no staleness.

### Index System

Indexes are first-class citizens in RaccoonDB, designed to work seamlessly with Tables:

- **IndexValueExtractor**: A function that extracts the indexable value from an object, defining what field(s) to index on.
- **Bucket Organization**: Index entries are grouped into buckets based on extracted values. Each unique value creates a bucket containing all object keys with that value.
- **Automatic Synchronization**: When objects are created, updated, or deleted via a Table, all attached indexes are automatically updated. Old index entries are removed and new ones added atomically.
- **Deferred Population**: Newly attached indexes don't auto-populate with existing objects. The `UpdateIndexes` method can be used to rebuild indexes from existing data.

## Design Patterns

### Prefix Namespacing
RaccoonDB uses key prefixes extensively to organize data logically within a single underlying store. Tables use prefixes like `objs/` for objects and `idxs/` for index data. This pattern enables multi-tenant separation and clean data organization without database overhead.

### Functional Types
The library provides functional programming primitives:
- **Option[T]**: Represents nullable values explicitly, avoiding nil pointer issues
- **Pair[T, U]**: Tuple type for combining related values

### Generic Type Safety
Go generics are used throughout to provide type-safe storage while abstracting serialization. `KeyObjectStore[T]` and `Table[T]` ensure compile-time type checking for stored objects.

## Typical Usage Flow

1. Create a base `KVStore` from an adapter (CoreKV, CometBFT, or Cosmos)
2. Optionally wrap with `PrefixStore` for logical namespacing
3. Create a `Table[T]` with an appropriate marshaler
4. Define and attach `Index` objects with extractors and marshalers
5. Use the table for CRUD operations - indexes stay synchronized automatically
