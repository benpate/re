package re

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {

	// Create a sample request
	body := bytes.NewReader([]byte("Hello World"))
	request, err := http.NewRequest("GET", "https://test.com", body)
	require.Nil(t, err)

	{
		first, err := ReadBody(request)
		require.Nil(t, err)
		require.Equal(t, []byte("Hello World"), first)
	}

	{
		second, err := ReadBody(request)
		require.Nil(t, err)
		require.Equal(t, "Hello World", string(second))
	}

	{
		third, err := io.ReadAll(request.Body)
		require.Nil(t, err)
		require.Equal(t, "Hello World", string(third))
	}
}

func TestRead_Error(t *testing.T) {

	// Create a sample request
	request, err := http.NewRequest("GET", "https://test.com", errorReader{})
	require.Nil(t, err)

	_, err = ReadBody(request)
	require.Error(t, err)
}
