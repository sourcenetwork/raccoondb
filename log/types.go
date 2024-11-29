package log

// TODO
// type Logger interface{}

// I could have a package level logger just to make things easier

// each store instance could optionally receive a log params / context thing which it can attenuate
// I could implement that with immutable trees, similar to how ctxs work
// that would allow for a span like messaging system.

// define an interface at the store level which lets me set the log params for each store
