package re

import (
	"bytes"
	"io"
	"net/http"

	"github.com/benpate/derp"
)

// ReadRequestBody reads the response.Body then replaces
// it with a new reader that can be read again by another process
func ReadRequestBody(request *http.Request) ([]byte, error) {

	// Request cannot be nil
	if request == nil {
		return []byte{}, nil
	}

	// If request.Body is nil, then return an empty byte array
	if request.Body == nil {
		return []byte{}, nil
	}

	originalBytes, err := io.ReadAll(request.Body)

	if err != nil {
		return []byte{}, derp.Wrap(err, "re.ReadRequestBody", "Error reading request body", derp.WithInternalError())
	}

	request.Body = io.NopCloser(bytes.NewReader(originalBytes))

	return originalBytes, nil
}

// ReadResponseBody reads the response.Body then replaces
// it with a new reader that can be read again by another process
//
// This is inspired by several articles, including:
// https://blog.manugarri.com/how-to-reuse-http-response-body-in-golang-2/
// https://medium.com/@xoen/golang-read-from-an-io-readwriter-without-loosing-its-content-2c6911805361
func ReadResponseBody(response *http.Response) ([]byte, error) {

	originalBytes, err := io.ReadAll(response.Body)

	if err != nil {
		return []byte{}, derp.Wrap(err, "re.ReadResponseBody", "Error reading response body", derp.WithInternalError())
	}

	response.Body = io.NopCloser(bytes.NewReader(originalBytes))

	return originalBytes, nil
}

// CloneResponse makes an exact copy of a response WITHOUT closing
// the original response body.
func CloneResponse(original *http.Response) http.Response {
	clone := *original
	originalBytes, _ := ReadResponseBody(original)
	clone.Body = io.NopCloser(bytes.NewReader(originalBytes))
	return clone
}
