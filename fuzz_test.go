package re

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

// FuzzReaderRoundTrip confirms that a Reader faithfully preserves arbitrary
// byte content through its accessor, chunked Read, and re-read paths, and never
// panics.
func FuzzReaderRoundTrip(f *testing.F) {

	f.Add([]byte(""))
	f.Add([]byte("Hello World"))
	f.Add([]byte{0x00, 0xFF, 0x10})
	f.Add(bytes.Repeat([]byte("x"), 1024)) // larger than io.ReadAll's initial buffer

	f.Fuzz(func(t *testing.T, data []byte) {

		// NewReaderFromBytes: accessors must return the exact input. (bytes.Equal
		// is used throughout so a nil and an empty slice compare equal, as they
		// should for this content.)
		fromBytes := NewReaderFromBytes(data)
		require.True(t, bytes.Equal(data, fromBytes.Bytes()), "Bytes() mismatch")
		require.Equal(t, string(data), fromBytes.String(), "String() mismatch")

		// NewReader (from an io.Reader) must capture the same content.
		fromReader, err := NewReader(bytes.NewReader(data))
		require.NoError(t, err)
		require.True(t, bytes.Equal(data, fromReader.Bytes()), "NewReader content mismatch")

		// Draining through a small buffer must reassemble the exact content, in
		// order — the io.Reader contract a cursor-less Read breaks.
		require.True(t, bytes.Equal(data, drainInChunks(t, fromBytes, 7)), "chunked Read mismatch")

		// After EOF the cursor rewinds, so the Reader can be drained again.
		again, err := io.ReadAll(fromBytes)
		require.NoError(t, err)
		require.True(t, bytes.Equal(data, again), "re-read content mismatch")

		// Close is always a no-op success.
		require.NoError(t, fromBytes.Close())
	})
}

// drainInChunks reads r to EOF through a buffer of the given size, returning the
// concatenation of every chunk. It fails the test on any non-EOF error.
func drainInChunks(t *testing.T, r io.Reader, bufSize int) []byte {
	t.Helper()

	var collected []byte

	// The buffer is allocated once in the loop's init clause and reused across
	// iterations (a fresh allocation per iteration would be wasteful).
	for p := make([]byte, bufSize); ; {
		n, err := r.Read(p)
		collected = append(collected, p[:n]...)

		if err == io.EOF {
			return collected
		}

		require.NoError(t, err)
	}
}
