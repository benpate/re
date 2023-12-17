package re

import (
	"io"

	"github.com/benpate/derp"
)

// Reader is a simple re-usable io.Reader.  It stores the entire contents of
// another io.Reader in memory, so it should not be used with large files.
//
// It implements the io.ReadCloser interface, along with a few other
// convenience methods.
//
// This is inspired by several articles, including:
// https://blog.flexicondev.com/read-go-http-request-body-multiple-times
type Reader struct {
	buffer []byte
}

// NewReader creates a new Reader from the given io.Reader
func NewReader(input io.Reader) (Reader, error) {
	var err error

	result := Reader{}

	result.buffer, err = io.ReadAll(input)

	if err != nil {
		return result, derp.Wrap(err, "re.NewReader", "Error reading input", derp.WithInternalError())
	}

	return result, nil
}

// NewReaderFromBytes creates a new Reader from the given slice of bytes
func NewReaderFromBytes(bytes []byte) Reader {
	return Reader{buffer: bytes}
}

// Read implements the io.Reader interface
func (r Reader) Read(p []byte) (int, error) {
	return copy(p, r.buffer), io.EOF
}

// Close implements the io.Closer interface
func (r Reader) Close() error {
	return nil
}

// Bytes returns the Reader's contents as a slice of bytes
func (r Reader) Bytes() []byte {
	return r.buffer
}

// String returns the Reader's contents as a string.
func (r Reader) String() string {
	return string(r.buffer)
}
