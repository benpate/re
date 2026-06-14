package re

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewReader_PreservesContent(t *testing.T) {

	reader, err := NewReader(strings.NewReader("some content"))
	require.NoError(t, err)
	require.Equal(t, "some content", reader.String())
	require.Equal(t, []byte("some content"), reader.Bytes())
}

func TestReader_ReadIntoFullBuffer(t *testing.T) {

	reader := NewReaderFromBytes([]byte("Hello World"))

	// A buffer at least as large as the content reads everything in one call
	p := make([]byte, 11)
	n, err := reader.Read(p)

	require.Equal(t, 11, n)
	require.Equal(t, "Hello World", string(p[:n]))
	require.ErrorIs(t, err, io.EOF)
}

func TestReader_EmptyInput(t *testing.T) {

	reader := NewReaderFromBytes([]byte{})
	require.Equal(t, "", reader.String())
	require.Equal(t, []byte{}, reader.Bytes())
}

// TestReader_ReadTruncatesLargeInput documents a known limitation: Reader.Read
// always returns io.EOF, even when it only copied part of the buffer. As a
// result, consumers like io.ReadAll that grow their buffer from a small initial
// size only receive the first chunk of a large input. This test pins that
// behavior so a future change (a fix) is noticed.
func TestReader_ReadTruncatesLargeInput(t *testing.T) {

	big := make([]byte, 5000)
	reader := NewReaderFromBytes(big)

	// A small buffer only receives len(p) bytes before EOF
	p := make([]byte, 512)
	n, err := reader.Read(p)

	require.Equal(t, 512, n)
	require.ErrorIs(t, err, io.EOF)

	// The full content is still available through the Bytes() accessor, which
	// does not go through the truncating Read path.
	require.Equal(t, 5000, len(reader.Bytes()))
}
