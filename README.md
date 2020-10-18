## [WIP] AIVEN OSS

Reads the repos of Aiven.io Github members and the work they have contributed to Open Source Software.

### The Ugly Way

You can use the magic of `WASM` to output these stats to a hideos HTML file. There is no caching, no lazy loading, nada. This was a learning exercise mostly.

### The `Go` Way

Run `aiven_oss.go`. At some point these will be separated out.

## Limitations

Go does not support selective imports based on `GOOS` or `GOARCH` and build tags do not support WASM (as far as I know)

## References

https://golangbot.com/webassembly-using-go/
https://www.aaron-powell.com/posts/2019-02-06-golang-wasm-3-interacting-with-js-from-go/