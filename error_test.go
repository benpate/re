package re

import "github.com/benpate/derp"

// errorReader is an io.Reader that always fails, used across the read and reader
// tests to exercise the error paths of anything that consumes a reader.
type errorReader struct{}

// Read always fails, and never returns any content.
func (e errorReader) Read(_ []byte) (n int, err error) {
	return 0, derp.Internal("errorReader", "ErrorReader always fails")
}
