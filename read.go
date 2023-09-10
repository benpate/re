package re

import (
	"net/http"

	"github.com/benpate/derp"
)

// ReadBody replaces the request.Body with a re.Reader
// and returns the body as a slice of bytes.
//
// This is inspired by several articles, including:
// https://blog.manugarri.com/how-to-reuse-http-response-body-in-golang-2/
// https://medium.com/@xoen/golang-read-from-an-io-readwriter-without-loosing-its-content-2c6911805361
func ReadBody(request *http.Request) ([]byte, error) {

	// If we already have a re.Reader, then read away.
	if reader, ok := request.Body.(Reader); ok {
		return reader.Bytes(), nil
	}

	// Otherwise, shim the re.Reader into the request.Body
	newReader, err := NewReader(request.Body)

	if err != nil {
		return nil, derp.Wrap(err, "re.ReadBody", "Error reading request body")
	}

	request.Body = newReader

	// Then return the body as a slice of bytes
	return newReader.Bytes(), nil
}
