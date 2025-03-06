.PHONY: test
test:
	go test ./...

.PHONY: test\:js
test\:js:
	GOOS=js GOARCH=wasm go test -exec="$$(go env GOROOT)/misc/wasm/go_js_wasm_exec" ./...

.PHONY: test\:wasip1
test\:wasip1:
	GOOS=wasip1 GOARCH=wasm go test -exec="$$(go env GOROOT)/misc/wasm/go_wasip1_wasm_exec" ./...
