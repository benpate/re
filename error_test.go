package re

import "github.com/benpate/derp"

type errorReader struct{}

func (e errorReader) Read(_ []byte) (n int, err error) {
	return 0, derp.Internal("errorReader", "ErrorReader always fails")
}
