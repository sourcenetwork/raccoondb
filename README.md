# Raccoon

Raccoon package is a general purpose graph store.
It wraps a KV store with a useful abstraction to work with graph structured data.

## Edges and Nodes

The main abstraction behind Raccoon is that of an `Edge`.
An Edge contains a Source Node and a Target Node.
The `Edge` interface specifies methods to return the target and source Nodes.

An Edge is a user defined type which gets serialized and persisted in the backend store.
Serialization is done by a `Marsheler`, an interface that marshals and unmarshals an Edge implementation.
Raccoon provides two default Marshalers, one for protobuff serialization and another for Json.

Users are free to add custom additional fields in their edge types.
Note that each Edge fully serializes its node data as well.


## Secondary Indexes

A notable fetaure of Raccoon is a lightweight secondary indexing feature.

Users may define secondary indexes in their schemas, which can be used for faster lookup.

Secondary indexes are built from Edge's through a `Mapper`.
A Mapper is a function that maps an Edge into a byte slice, which is used as the indexing key.

## Keying System

Edge nodes must map to unique keys which are used internally by Raccoon.
Edges are unique identified by the keys of their source and target nodes.

Key generation is specified through the `NodeKeyer` interface.
Users must implement a NodeKeyer for their custom Node types.

## Example usage

```go
// TODO
```
