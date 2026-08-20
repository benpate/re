# re ♻️

[![Go Reference](https://pkg.go.dev/badge/github.com/benpate/re.svg)](https://pkg.go.dev/github.com/benpate/re)
[![Version](https://img.shields.io/github/v/release/benpate/re?include_prereleases&style=flat-square&color=brightgreen)](https://github.com/benpate/re/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/benpate/re/go.yml?branch=main)](https://github.com/benpate/re/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/re?style=flat-square)](https://goreportcard.com/report/github.com/benpate/re)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/re.svg?style=flat-square)](https://codecov.io/gh/benpate/re)

## Re-Reader for Go

In some circumstances, it's necessary to use a Golang [io.Reader](https://pkg.go.dev/io#Reader) multiple times. While readers themselves don't support this, we can fake it by copying the value into a buffer and then replaying the values whenever asked.

This is useful for processing HTTP requests, particularly when multiple layers of an app all need the same access to the HTTP request body. This happens, for example, when an [http signature](https://github.com/benpate/hannibal/tree/main/sigs) middleware needs to validate the body digest before allowing the transaction to pass through to the main handler app.

### A Note on Memory and Limits

Because `re` keeps the whole dataset in memory and replays it, it trades memory for re-readability. Don't reach for it with very large values or long-running streaming readers — only when the data is small and you reasonably expect to read it a second time.

To keep that tradeoff from becoming a liability, the HTTP body helpers take a `maxBytes` limit. A body larger than the limit returns an error instead of allocating without bound or silently parsing a truncated payload, so a hostile peer can't force a runaway allocation. Pass `re.DefaultMaximum` (1 MiB — comfortable for typical JSON payloads) when you have no stricter requirement of your own.

## Read from an io.Reader Multiple Times

Wrap any reader in a `re.Reader`. This copies all of its bytes into a memory buffer that can be re-read as many times as necessary. Each `Read` cycle ends with `io.EOF` and rewinds automatically, so the next consumer starts from the beginning.

```go
func main() {
	// A single-use reader
	singleReader := strings.NewReader("Hello World.")

	// Wrap it so it can be read multiple times
	multipleReader, err := re.NewReader(singleReader)
	if err != nil {
		panic(err)
	}

	readAndPrint(multipleReader) // Hello World.
	readAndPrint(multipleReader) // Hello World.
	readAndPrint(multipleReader) // Hello World.
}

func readAndPrint(r io.Reader) {
	b, _ := io.ReadAll(r)
	fmt.Println(string(b))
}
```

A `re.Reader` is also useful on its own, away from `io.Reader` plumbing:

- `re.NewReaderFromBytes(b)` wraps a slice you already have (the slice is retained, not copied).
- `Reset()` rewinds the cursor early, without draining to `io.EOF` first.
- `Bytes()` and `String()` return the full contents directly, independent of the read cursor.

A `Reader` is **not** safe for concurrent use — the cursor is shared mutable state. To fan out to parallel consumers, give each its own reader over the same bytes with `re.NewReaderFromBytes(r.Bytes())`.

## Read the Body from an http.Request or http.Response

This library also includes helper functions to read from [`http.Request`](https://pkg.go.dev/net/http#Request) and [`http.Response`](https://pkg.go.dev/net/http#Response) bodies. These functions do not use a `re.Reader`, but simply read the existing Body reader (up to `maxBytes`) then replace it with a fresh reader that can be read again by another process. They guard nil requests, responses, and bodies, so a nil body (legal on a HEAD response) yields an empty slice instead of a panic.

```go
func handler(request *http.Request) {
	body, err := re.ReadRequestBody(request, re.DefaultMaximum)
}
```

`re.ReadResponseBody(response, maxBytes)` does the same for a response. And `re.CloneResponse(response, maxBytes)` makes an independent copy of a response — reading and restoring the original's body so both the original and the clone remain readable. The clone gets its own headers and trailers too, so you can modify the copy without touching the original.

## Issues and Pull Requests Welcome

As is everything in life, `re` is a work in progress and will benefit from your experience reports, use cases, and contributions. If you have an idea for making this library better, send in a pull request. We're all in this together! ♻️
