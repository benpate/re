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
		first, err := ReadRequestBody(request, DefaultMaximum)
		require.Nil(t, err)
		require.Equal(t, []byte("Hello World"), first)
	}

	{
		second, err := ReadRequestBody(request, DefaultMaximum)
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
		first, err := ReadResponseBody(&response, DefaultMaximum)
		require.Nil(t, err)
		require.Equal(t, []byte("Hello World"), first)
	}

	{
		second, err := ReadResponseBody(&response, DefaultMaximum)
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

	result, err := ReadRequestBody(request, DefaultMaximum)
	require.Error(t, err)
	require.Equal(t, []byte{}, result)
}

func TestReadResponse_Error(t *testing.T) {

	// Create a sample request
	response := http.Response{
		Body: io.NopCloser(errorReader{}),
	}

	result, err := ReadResponseBody(&response, DefaultMaximum)
	require.Error(t, err)
	require.Equal(t, []byte{}, result)
}

func TestReadRequestBody_ExceedsMaximum(t *testing.T) {

	// A body larger than the limit must error and return no bytes (never a truncated body)
	request, err := http.NewRequest("POST", "https://test.com", bytes.NewReader([]byte("Hello World")))
	require.Nil(t, err)

	result, err := ReadRequestBody(request, 5)
	require.Error(t, err)
	require.Equal(t, []byte{}, result)
}

func TestReadRequestBody_AtMaximum(t *testing.T) {

	// A body exactly at the limit is read successfully
	request, err := http.NewRequest("POST", "https://test.com", bytes.NewReader([]byte("Hello")))
	require.Nil(t, err)

	result, err := ReadRequestBody(request, 5)
	require.NoError(t, err)
	require.Equal(t, "Hello", string(result))
}

func TestReadResponseBody_ExceedsMaximum(t *testing.T) {

	// A response body larger than the limit must error and return no bytes
	response := http.Response{Body: io.NopCloser(bytes.NewReader([]byte("Hello World")))}

	result, err := ReadResponseBody(&response, 5)
	require.Error(t, err)
	require.Equal(t, []byte{}, result)
}

func TestCloneResponse_ExceedsMaximum(t *testing.T) {

	// Cloning a response whose body exceeds the limit must propagate the error
	original := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte("Hello World"))),
	}

	_, err := CloneResponse(original, 5)
	require.Error(t, err)
}
