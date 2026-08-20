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
	result, err := ReadRequestBody(nil, DefaultMaximum)
	require.NoError(t, err)
	require.Equal(t, []byte{}, result)
}

func TestReadRequestBody_NilBody(t *testing.T) {

	// A request with no body yields an empty body and no error
	request, err := http.NewRequest("GET", "https://test.com", nil)
	require.NoError(t, err)
	require.Nil(t, request.Body)

	result, err := ReadRequestBody(request, DefaultMaximum)
	require.NoError(t, err)
	require.Equal(t, []byte{}, result)
}

func TestReadResponseBody_NilResponse(t *testing.T) {

	// A nil response yields an empty body and no error (no panic)
	result, err := ReadResponseBody(nil, DefaultMaximum)
	require.NoError(t, err)
	require.Equal(t, []byte{}, result)
}

func TestReadResponseBody_NilBody(t *testing.T) {

	// A response with no body yields an empty body and no error (no panic).
	// A nil Body is legal, e.g. on a response to a HEAD request.
	response := &http.Response{}

	result, err := ReadResponseBody(response, DefaultMaximum)
	require.NoError(t, err)
	require.Equal(t, []byte{}, result)
}

func TestCloneResponse_NilOriginal(t *testing.T) {

	// A nil original yields a zero-value response and no error, without panicking
	require.NotPanics(t, func() {
		clone, err := CloneResponse(nil, DefaultMaximum)
		require.NoError(t, err)
		require.Equal(t, http.Response{}, clone)
	})
}

func TestCloneResponse_NilBody(t *testing.T) {

	// Cloning a response with a nil body must not panic
	require.NotPanics(t, func() {
		clone, err := CloneResponse(&http.Response{StatusCode: 204}, DefaultMaximum)
		require.NoError(t, err)
		require.Equal(t, 204, clone.StatusCode)
	})
}

func TestCloneResponse_ReadError(t *testing.T) {

	// When the original body cannot be read, the error is propagated rather
	// than silently producing an empty clone.
	original := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(errorReader{}),
	}

	_, err := CloneResponse(original, DefaultMaximum)
	require.Error(t, err)
}

func TestCloneResponse(t *testing.T) {

	original := &http.Response{
		StatusCode: 200,
		Header:     http.Header{"X-Test": []string{"value"}},
		Body:       io.NopCloser(bytes.NewReader([]byte("Hello World"))),
	}

	clone, err := CloneResponse(original, DefaultMaximum)
	require.NoError(t, err)

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

	clone, err := CloneResponse(original, DefaultMaximum)
	require.NoError(t, err)

	// Reading the clone to completion must not affect the original's reader
	_, err = io.ReadAll(clone.Body)
	require.NoError(t, err)

	originalBody, err := io.ReadAll(original.Body)
	require.NoError(t, err)
	require.Equal(t, "payload", string(originalBody))
}

func TestCloneResponse_IndependentHeaders(t *testing.T) {

	// A clone must own its headers outright. A caller who copies a response in
	// order to change it would otherwise rewrite the original's headers too.
	original := &http.Response{
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{"text/plain"}},
		Body:       io.NopCloser(bytes.NewReader([]byte("Hello World"))),
	}

	clone, err := CloneResponse(original, DefaultMaximum)
	require.NoError(t, err)

	clone.Header.Set("Content-Type", "application/json")
	clone.Header.Set("X-Added", "yes")

	require.Equal(t, "text/plain", original.Header.Get("Content-Type"))
	require.Equal(t, "", original.Header.Get("X-Added"))
}

func TestCloneResponse_IndependentTrailers(t *testing.T) {

	original := &http.Response{
		StatusCode: 200,
		Trailer:    http.Header{"Expires": []string{"never"}},
		Body:       io.NopCloser(bytes.NewReader([]byte("Hello World"))),
	}

	clone, err := CloneResponse(original, DefaultMaximum)
	require.NoError(t, err)

	clone.Trailer.Set("Expires", "tomorrow")

	require.Equal(t, "never", original.Trailer.Get("Expires"))
}

func TestCloneResponse_IndependentTransferEncoding(t *testing.T) {

	original := &http.Response{
		StatusCode:       200,
		TransferEncoding: []string{"chunked"},
		Body:             io.NopCloser(bytes.NewReader([]byte("Hello World"))),
	}

	clone, err := CloneResponse(original, DefaultMaximum)
	require.NoError(t, err)
	require.Equal(t, []string{"chunked"}, clone.TransferEncoding)

	clone.TransferEncoding[0] = "identity"

	require.Equal(t, []string{"chunked"}, original.TransferEncoding)
}

func TestCloneResponse_NilHeadersStayNil(t *testing.T) {

	// Copying must not invent maps that the original did not have, since a nil
	// Header and an empty one are not the same thing to net/http.
	original := &http.Response{
		StatusCode: 204,
		Body:       io.NopCloser(bytes.NewReader(nil)),
	}

	clone, err := CloneResponse(original, DefaultMaximum)
	require.NoError(t, err)

	require.Nil(t, clone.Header)
	require.Nil(t, clone.Trailer)
	require.Nil(t, clone.TransferEncoding)
}
