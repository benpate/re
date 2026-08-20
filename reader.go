package re

import (
	"io"

	"github.com/benpate/derp"
)

// Reader is a simple, re-usable io.ReadCloser. It stores the entire contents of
// another io.Reader in memory, so it should not be used with large files.
//
// Unlike a one-shot reader, a Reader can be read repeatedly: each Read advances
// an internal cursor and returns io.EOF when the content is exhausted, then the
// cursor rewinds so the next Read cycle starts again from the beginning. This
// lets one consumer read the content, then a later consumer read it again.
//
// A Reader is NOT safe for concurrent use: the cursor is shared mutable state,
// so the re-reads must happen in sequence, not in parallel. To fan out to
// concurrent consumers, give each its own Reader over the same bytes (for
// example, NewReaderFromBytes(r.Bytes())).
type Reader struct {
	buffer []byte // immutable backing content; accessors and re-reads draw from this
	pos    int    // cursor: index of the next unread byte (not safe for concurrent Read)
}

// NewReader creates a new Reader by reading the given io.Reader to completion.
func NewReader(input io.Reader) (*Reader, error) {

	const location = "re.NewReader"

	// The re-readable-buffer approach is inspired by several articles, including:
	// https://blog.flexicondev.com/read-go-http-request-body-multiple-times
	buffer, err := io.ReadAll(input)

	if err != nil {
		return &Reader{}, derp.Wrap(err, location, "Unable to read input", derp.WithInternalError())
	}

	return &Reader{buffer: buffer}, nil
}

// NewReaderFromBytes creates a new Reader from the given slice of bytes. The
// slice is retained, not copied, so callers must not mutate it afterward.
func NewReaderFromBytes(buffer []byte) *Reader {
	return &Reader{buffer: buffer}
}

// Read implements the io.Reader interface. It copies the next unread bytes into
// p, advancing the cursor. When the content is exhausted it returns 0, io.EOF
// and rewinds the cursor, so the Reader can be read again from the start.
func (r *Reader) Read(p []byte) (int, error) {

	// RULE: At end of content, signal EOF and rewind for the next read cycle.
	// The rewind is what makes this Reader re-readable.
	if r.pos >= len(r.buffer) {
		r.pos = 0
		return 0, io.EOF
	}

	// Otherwise fill as much of p as the remaining content allows, and advance.
	n := copy(p, r.buffer[r.pos:])
	r.pos += n

	// Another satisfied customer
	return n, nil
}

// Reset rewinds the cursor to the start of the content, so the next Read returns
// the content from the beginning.
func (r *Reader) Reset() {
	r.pos = 0
}

// Close implements the io.Closer interface. It is a no-op (the content is held
// in memory) and never returns an error.
func (r *Reader) Close() error {

	// Nothing to close, nothing to see here
	return nil
}

// Bytes returns the Reader's full contents, independent of the read cursor.
func (r *Reader) Bytes() []byte {
	return r.buffer
}

// String returns the Reader's full contents as a string, independent of the
// read cursor.
func (r *Reader) String() string {
	return string(r.buffer)
}
