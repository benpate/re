package re

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReader_ReReader(t *testing.T) {

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
