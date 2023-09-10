package re

import "github.com/benpate/derp"

type errorReader struct{}

func (e errorReader) Read(p []byte) (n int, err error) {
	return 0, derp.NewInternalError("errorReader", "ErrorReader always fails")
}
