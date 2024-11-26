package store

// base kv - maybe impl a general purpose kv wrapper which handles basic validation
// and consistent errors for any kv implementation?
// that way validations could be concentrated here
// eg. iterator validations, set / get key / value validation and so on.

// maybe move the definition of the kv to the stores package
// move definition of iterator to the iterator pkg

// sort the errors
