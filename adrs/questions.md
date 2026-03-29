# Questions for Future ADRs

## History and Motivation
- What application or use case drove creating RaccoonDB?
- What problems did existing KV store libraries not solve?

## Versioning (v2 Branch)
- What changed from v1 to v2?
- Are there breaking changes in the interface hierarchy?
- Should v1 still be maintained or is it deprecated?

## Concurrency Model
- Is the library expected to be thread-safe?
- Are individual stores safe for concurrent read/write access?
- Is concurrency safety the caller's responsibility or handled internally?
- Do different adapters have different concurrency guarantees?

## Relation Primitive
- What is the purpose of `primitives/relation.go`?
- Is it work-in-progress or a complete feature?
- How does it fit into the primitives hierarchy?

## Testing Conventions
- What is the standard pattern for validating new store adapters?
- How should new primitives be tested?
- Is `store/test` meant to be a compliance test suite for adapters?

## Known Limitations and Future Work
- Are there known pain points or areas that need improvement?
- Any planned features or primitives?
- Performance considerations or known bottlenecks?
