# RaccoonDB

RaccoonDB is a library which builds on top of a Key-Value Store interface and provides higher level constructs such as indexes, counters and tables.
The purpose of Raccoon is to bootstrap developing data management applications on top of arbitrary KV Stores.

Raccoon's design leverages the well-known Iterator interface as the primary abstraction over sequences.
Iterators can be composed to add extra functionality such as mapping, filtering, merging and so on.