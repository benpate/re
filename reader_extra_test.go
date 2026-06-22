package re

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestReader_SequentialConsumers documents the supported usage: one consumer
// reads the content to EOF, then a later consumer reads the same content again.
// (Concurrent reads of one Reader are not supported; each consumer that needs
// parallel access should get its own Reader over the same bytes.)
func TestReader_SequentialConsumers(t *testing.T) {

	shared := NewReaderFromBytes([]byte("shared body"))

	for consumer := 0; consumer < 3; consumer++ {
		got, err := io.ReadAll(shared)
		require.NoError(t, err)
		require.Equal(t, "shared body", string(got))
	}
}

func TestNewReader_PreservesContent(t *testing.T) {

	reader, err := NewReader(strings.NewReader("some content"))
	require.NoError(t, err)
	require.Equal(t, "some content", reader.String())
	require.Equal(t, []byte("some content"), reader.Bytes())
}

func TestReader_ReadIntoFullBuffer(t *testing.T) {

	reader := NewReaderFromBytes([]byte("Hello World"))

	// A buffer at least as large as the content reads everything in one call.
	// Per the io.Reader contract, that call returns the data with a nil error;
	// EOF arrives on the following call.
	p := make([]byte, 11)
	n, err := reader.Read(p)

	require.Equal(t, 11, n)
	require.Equal(t, "Hello World", string(p[:n]))
	require.NoError(t, err)

	// The next Read reports EOF (and rewinds for re-reading).
	n, err = reader.Read(p)
	require.Equal(t, 0, n)
	require.ErrorIs(t, err, io.EOF)
}

func TestReader_EmptyInput(t *testing.T) {

	reader := NewReaderFromBytes([]byte{})
	require.Equal(t, "", reader.String())
	require.Equal(t, []byte{}, reader.Bytes())
}

// TestReader_ReadLargeInputViaSmallBuffer confirms that a large input is read in
// full, in order, when consumed through a buffer smaller than the content. This
// is the io.Reader contract that a cursor-less implementation would violate
// (truncating to the first chunk and repeating it).
func TestReader_ReadLargeInputViaSmallBuffer(t *testing.T) {

	// Distinct, position-dependent content so any repetition or gap is detectable.
	big := make([]byte, 5000)
	for i := range big {
		big[i] = byte(i % 251) // 251 is prime, so the pattern doesn't align to the buffer size
	}

	reader := NewReaderFromBytes(big)

	// Drain through a 512-byte buffer; the full content must come back in order.
	collected := drainInChunks(t, reader, 512)
	require.Equal(t, big, collected, "small-buffer reads must reassemble the full content in order")
}

// TestReader_ReadAllLargeInput confirms io.ReadAll receives the entire content
// even when it exceeds io.ReadAll's initial buffer size (512 bytes) — the exact
// case the previous cursor-less implementation truncated.
func TestReader_ReadAllLargeInput(t *testing.T) {

	big := make([]byte, 5000)
	for i := range big {
		big[i] = byte(i % 251)
	}

	got, err := io.ReadAll(NewReaderFromBytes(big))

	require.NoError(t, err)
	require.Equal(t, big, got)
}

// TestReader_RereadAfterEOF confirms the cursor rewinds on EOF, so the same
// Reader can be drained repeatedly — the type's core purpose.
func TestReader_RereadAfterEOF(t *testing.T) {

	reader := NewReaderFromBytes([]byte("Hello World"))

	for i := 0; i < 3; i++ {
		got, err := io.ReadAll(reader)
		require.NoError(t, err)
		require.Equal(t, "Hello World", string(got))
	}
}

// TestReader_Reset confirms Reset rewinds the cursor mid-stream.
func TestReader_Reset(t *testing.T) {

	reader := NewReaderFromBytes([]byte("Hello World"))

	p := make([]byte, 5)
	_, err := reader.Read(p) // consume "Hello"
	require.NoError(t, err)
	require.Equal(t, "Hello", string(p))

	reader.Reset()

	got, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, "Hello World", string(got), "Reset must rewind to the start")
}
