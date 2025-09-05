# re 🔖

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/benpate/re)
[![Version](https://img.shields.io/github/v/release/benpate/re?include_prereleases&style=flat-square&color=brightgreen)](https://github.com/benpate/re/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/benpate/re/go.yml?branch=main)](https://github.com/benpate/re/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/re?style=flat-square)](https://goreportcard.com/report/github.com/benpate/re)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/re.svg?style=flat-square)](https://codecov.io/gh/benpate/re)

## Re-Reader for Go
In some circumstances, it's necessary to us a Golang [io.Reader](https://pkg.go.dev/io#Reader) multiple times.  While readers themselves don't support this, we can fake it by copying the value into a buffer and then replaying the values whenever asked.

**WARNING:** This library uses more resources than a regular reader.  In particular, it's a bad idea to use this for very large values (because it keeps the whole dataset in memory) or long running streaming readers (because it reads and replays the entire value).  It should only be used when the data is small and you reasonably believe you will need to use the reader a second time.

However, it is useful for processing HTTP requests, particularly when multiple layers of an app all need the same access to the HTTP request body.  This happens, for example, when an [http signature](https://github.com/benpate/hannibal/tree/main/sigs) middleware needs to validate the body digest before allowing the transaction to pass through to the main handler app.

## Read from an io.Reader Multiple Times

To use this library, just wrap any reader in a `re.Reader`.  This will copy all of its values into a memory buffer which can be re-read as many times as necessary.

```go
func main() {
	// First, a single-use reader 
	singleReader := strings.NewReader(("Hello World."))

	// This reader will work multiple times
	multipleReader := re.NewReader(singleReader)

    readAndPrint(multipleReader)
    readAndPrint(multipleReader)
    readAndPrint(multipleReader)
}

func readAndPrint(r io.Reader) {
    b, _ := io.ReadAll(r)
    fmt.Println(string(b))
}
	
```

## Read the Body from an http.Request or http.Response

This library also includes helper functions to read from [`http.Request`](https://pkg.go.dev/net/http#Request) and [`http.Response`](https://pkg.go.dev/net/http#Response) bodies. These functions do not use a re.Reader, but simply read the existing Body reader then replace it with a fresh reader that can be read again by another process.


```go
func handler(request *http.Request) {
	body, err := re.ReadRequestBody(request)
}
```

## Issues and Pull Requests Welcome

As is everything in life, `re` is a work in progress and will benefit from your experience reports, use cases, and contributions.  If you have an idea for making this library better, send in a pull request.  We're all in this together! 🔖
