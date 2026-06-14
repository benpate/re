package re

import (
	"bytes"
	"testing"
)

// FuzzReaderRoundTrip confirms that a Reader faithfully preserves arbitrary
// byte content through its accessor and full-buffer Read paths, and never
// panics.
func FuzzReaderRoundTrip(f *testing.F) {

	f.Add([]byte(""))
	f.Add([]byte("Hello World"))
	f.Add([]byte{0x00, 0xFF, 0x10})
	f.Add(bytes.Repeat([]byte("x"), 1024)) // larger than io.ReadAll's initial buffer

	f.Fuzz(func(t *testing.T, data []byte) {

		// NewReaderFromBytes: accessors must return the exact input
		fromBytes := NewReaderFromBytes(data)
		if !bytes.Equal(fromBytes.Bytes(), data) {
			t.Fatalf("Bytes() mismatch: got %d bytes, want %d", len(fromBytes.Bytes()), len(data))
		}
		if fromBytes.String() != string(data) {
			t.Fatal("String() mismatch")
		}

		// NewReader (from an io.Reader) must capture the same content
		fromReader, err := NewReader(bytes.NewReader(data))
		if err != nil {
			t.Fatalf("NewReader error: %v", err)
		}
		if !bytes.Equal(fromReader.Bytes(), data) {
			t.Fatal("NewReader content mismatch")
		}

		// A single Read into a full-size buffer returns all of the content
		p := make([]byte, len(data))
		n, _ := fromBytes.Read(p)
		if n != len(data) {
			t.Fatalf("Read returned %d bytes, want %d", n, len(data))
		}
		if !bytes.Equal(p, data) {
			t.Fatal("Read content mismatch")
		}

		// Close is always a no-op success
		if err := fromBytes.Close(); err != nil {
			t.Fatalf("Close error: %v", err)
		}
	})
}
