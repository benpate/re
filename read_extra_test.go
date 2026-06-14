package re

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadRequestBody_NilRequest(t *testing.T) {

	// A nil request yields an empty body and no error
	result, err := ReadRequestBody(nil)
	require.NoError(t, err)
	require.Equal(t, []byte{}, result)
}

func TestReadRequestBody_NilBody(t *testing.T) {

	// A request with no body yields an empty body and no error
	request, err := http.NewRequest("GET", "https://test.com", nil)
	require.NoError(t, err)
	require.Nil(t, request.Body)

	result, err := ReadRequestBody(request)
	require.NoError(t, err)
	require.Equal(t, []byte{}, result)
}

func TestCloneResponse(t *testing.T) {

	original := &http.Response{
		StatusCode: 200,
		Header:     http.Header{"X-Test": []string{"value"}},
		Body:       io.NopCloser(bytes.NewReader([]byte("Hello World"))),
	}

	clone := CloneResponse(original)

	// The clone carries the same metadata
	require.Equal(t, 200, clone.StatusCode)
	require.Equal(t, "value", clone.Header.Get("X-Test"))

	// The clone's body is readable and matches the original content
	cloneBody, err := io.ReadAll(clone.Body)
	require.NoError(t, err)
	require.Equal(t, "Hello World", string(cloneBody))

	// The ORIGINAL body is still readable (not consumed by the clone)
	originalBody, err := io.ReadAll(original.Body)
	require.NoError(t, err)
	require.Equal(t, "Hello World", string(originalBody))
}

func TestCloneResponse_IndependentBodies(t *testing.T) {

	original := &http.Response{
		StatusCode: 201,
		Body:       io.NopCloser(bytes.NewReader([]byte("payload"))),
	}

	clone := CloneResponse(original)

	// Reading the clone to completion must not affect the original's reader
	_, err := io.ReadAll(clone.Body)
	require.NoError(t, err)

	originalBody, err := io.ReadAll(original.Body)
	require.NoError(t, err)
	require.Equal(t, "payload", string(originalBody))
}
