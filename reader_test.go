package re

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReader_ReReader(t *testing.T) {

	// Create a sample request
	reader, err := NewReader(strings.NewReader("Hello World"))
	require.Nil(t, err)

	for i := 0; i < 100; i++ {
		first, err := io.ReadAll(reader)
		require.Nil(t, err)
		require.Equal(t, "Hello World", string(first))
	}
}

func TestReader_ReReaderBytes(t *testing.T) {

	// Create a sample request
	reader := NewReaderFromBytes([]byte("Hello World"))

	for i := 0; i < 100; i++ {
		first, err := io.ReadAll(reader)
		require.Nil(t, err)
		require.Equal(t, "Hello World", string(first))
	}
}

func TestReader_Helpers(t *testing.T) {

	// Create a sample request
	reader := NewReaderFromBytes([]byte("Hello World"))

	require.Equal(t, "Hello World", reader.String())
	require.Equal(t, []byte("Hello World"), reader.Bytes())
	require.Nil(t, reader.Close())
}

func TestReader_Error(t *testing.T) {
	_, err := NewReader(errorReader{})
	require.NotNil(t, err)
}
