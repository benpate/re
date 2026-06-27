package re

import (
	"bytes"
	"io"
	"net/http"

	"github.com/benpate/derp"
)

// DefaultMaximum is a sane default body-size limit (1 MiB) for callers that have
// no stricter requirement of their own. It is large enough for typical JSON API
// payloads while preventing a hostile peer from forcing an unbounded allocation.
const DefaultMaximum = int64(1 << 20)

// ReadRequestBody reads the request.Body (up to maxBytes) then replaces it with a
// new reader that can be read again by another process. If the body exceeds
// maxBytes, it returns an error and reads nothing further, so callers never parse
// a truncated body. Pass DefaultMaximum if you have no stricter limit.
func ReadRequestBody(request *http.Request, maxBytes int64) ([]byte, error) {

	// Request cannot be nil
	if request == nil {
		return []byte{}, nil
	}

	// If request.Body is nil, then return an empty byte array
	if request.Body == nil {
		return []byte{}, nil
	}

	// Cap the body so an oversized payload errors instead of allocating without bound.
	originalBytes, err := io.ReadAll(http.MaxBytesReader(nil, request.Body, maxBytes))

	if err != nil {
		return []byte{}, derp.Wrap(err, "re.ReadRequestBody", "Unable to read request body (or body exceeds maximum size)", maxBytes, derp.WithInternalError())
	}

	request.Body = io.NopCloser(bytes.NewReader(originalBytes))

	return originalBytes, nil
}

// ReadResponseBody reads the response.Body (up to maxBytes) then replaces it with
// a new reader that can be read again by another process. If the body exceeds
// maxBytes, it returns an error and reads nothing further, so callers never parse
// a truncated body. Pass DefaultMaximum if you have no stricter limit.
func ReadResponseBody(response *http.Response, maxBytes int64) ([]byte, error) {

	// This is inspired by several articles, including:
	// https://blog.manugarri.com/how-to-reuse-http-response-body-in-golang-2/
	// https://medium.com/@xoen/golang-read-from-an-io-readwriter-without-loosing-its-content-2c6911805361

	// Response cannot be nil
	if response == nil {
		return []byte{}, nil
	}

	// If response.Body is nil, then return an empty byte array. A nil Body is
	// legal (for instance, on a response to a HEAD request).
	if response.Body == nil {
		return []byte{}, nil
	}

	// Cap the body so an oversized payload errors instead of allocating without bound.
	originalBytes, err := io.ReadAll(http.MaxBytesReader(nil, response.Body, maxBytes))

	if err != nil {
		return []byte{}, derp.Wrap(err, "re.ReadResponseBody", "Unable to read response body (or body exceeds maximum size)", maxBytes, derp.WithInternalError())
	}

	response.Body = io.NopCloser(bytes.NewReader(originalBytes))

	return originalBytes, nil
}

// CloneResponse makes an exact copy of a response (reading up to maxBytes of its
// body) WITHOUT closing the original response body. A nil original yields a
// zero-value response. Pass DefaultMaximum if you have no stricter limit.
func CloneResponse(original *http.Response, maxBytes int64) (http.Response, error) {

	// Original cannot be nil (dereferencing it below would panic).
	if original == nil {
		return http.Response{}, nil
	}

	clone := *original

	originalBytes, err := ReadResponseBody(original, maxBytes)

	if err != nil {
		return clone, derp.Wrap(err, "re.CloneResponse", "Unable to read original response body")
	}

	clone.Body = io.NopCloser(bytes.NewReader(originalBytes))
	return clone, nil
}
