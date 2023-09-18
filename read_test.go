package re

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadRequest(t *testing.T) {

	// Create a sample request
	body := bytes.NewReader([]byte("Hello World"))
	request, err := http.NewRequest("GET", "https://test.com", body)
	require.Nil(t, err)

	{
		first, err := ReadRequestBody(request)
		require.Nil(t, err)
		require.Equal(t, []byte("Hello World"), first)
	}

	{
		second, err := ReadRequestBody(request)
		require.Nil(t, err)
		require.Equal(t, "Hello World", string(second))
	}

	{
		third, err := io.ReadAll(request.Body)
		require.Nil(t, err)
		require.Equal(t, "Hello World", string(third))
	}
}

func TestReadResponse(t *testing.T) {

	// Create a sample request
	body := bytes.NewReader([]byte("Hello World"))
	response := http.Response{}
	response.Body = io.NopCloser(body)

	{
		first, err := ReadResponseBody(&response)
		require.Nil(t, err)
		require.Equal(t, []byte("Hello World"), first)
	}

	{
		second, err := ReadResponseBody(&response)
		require.Nil(t, err)
		require.Equal(t, "Hello World", string(second))
	}

	{
		third, err := io.ReadAll(response.Body)
		require.Nil(t, err)
		require.Equal(t, "Hello World", string(third))
	}
}

func TestReadRequest_Error(t *testing.T) {

	// Create a sample request
	request, err := http.NewRequest("GET", "https://test.com", errorReader{})
	require.Nil(t, err)

	result, err := ReadRequestBody(request)
	require.Error(t, err)
	require.Equal(t, []byte{}, result)
}

func TestReadResponse_Error(t *testing.T) {

	// Create a sample request
	response := http.Response{
		Body: io.NopCloser(errorReader{}),
	}

	result, err := ReadResponseBody(&response)
	require.Error(t, err)
	require.Equal(t, []byte{}, result)
}
