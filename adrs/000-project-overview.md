---
date: 2026-03-29
---

# RaccoonDB

RaccoonDB is a Go library which provides common primitives and abstractions for KV Stores, the goal is to simplify KV Store usage while building apps.

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

### Primitives 

Primitives are the building blocks for application developrs, meant to represent reusable and composable primitives over KVStores.

Some examples are: a counted kv store, byte indexes and object store (which handles marshaling).


### Table

The Table abstraction is a top level abstraction which composes several of the underlying primitives.
It's meant to model a data storage abstraction which has some of the niceties of a database table developers are used from the DB world.

Tables can have indices defined on them, synchronizes index remove and updating while creating object, asserting there is no staleness.
