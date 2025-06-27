package httphelpers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	helpers "github.com/launchdarkly/go-test-helpers/v3"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChunkedStreamingHandlerReturnsResponseBeforeFirstData(t *testing.T) {
	handler, stream := ChunkedStreamingHandler(nil, "text/plain")
	defer stream.Close()

	WithServer(handler, func(server *httptest.Server) {
		resp, err := http.DefaultClient.Get(server.URL)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestChunkedStreamingHandlerSend(t *testing.T) {
	initialData := []byte("hello,")
	handler, stream := ChunkedStreamingHandler(initialData, "text/plain")
	defer stream.Close()

	stream.Enqueue([]byte("first,"))
	stream.Send([]byte("this isn't sent because there are no connections"))
	stream.Enqueue([]byte("second,"))

	WithServer(handler, func(server *httptest.Server) {
		resp1, err := http.DefaultClient.Get(server.URL)
		require.NoError(t, err)
		defer resp1.Body.Close()

		assert.Equal(t, 200, resp1.StatusCode)
		assert.Equal(t, "text/plain", resp1.Header.Get("Content-Type"))

		stream.Send(nil)             // should have no effect
		stream.Send(make([]byte, 0)) // also no effect
		stream.Send([]byte("third,"))

		expected := "hello,first,second,third,"
		assert.Equal(t, expected, string(helpers.ReadWithTimeout(resp1.Body, len(expected), time.Second)))

		resp2, err := http.DefaultClient.Get(server.URL)
		require.NoError(t, err)
		defer resp2.Body.Close()

		expected = "hello,"
		assert.Equal(t, expected, string(helpers.ReadWithTimeout(resp2.Body, len(expected), time.Second)))

		stream.Send([]byte("fourth."))
		expected = "fourth."
		assert.Equal(t, expected, string(helpers.ReadWithTimeout(resp1.Body, len(expected), time.Second)))
		assert.Equal(t, expected, string(helpers.ReadWithTimeout(resp2.Body, len(expected), time.Second)))
	})
}

func TestChunkedStreamingHandlerEndAll(t *testing.T) {
	initialData := []byte("hello,")
	handler, stream := ChunkedStreamingHandler(initialData, "text/plain")
	defer stream.Close()

	WithServer(handler, func(server *httptest.Server) {
		resp1, err := http.DefaultClient.Get(server.URL)
		require.NoError(t, err)
		defer resp1.Body.Close()

		go func() {
			stream.Send([]byte("goodbye."))
			stream.EndAll()
		}()

		// ReadAll won't return until the stream is closed
		data, err := io.ReadAll(resp1.Body)
		require.NoError(t, err)
		assert.Equal(t, "hello,goodbye.", string(data))

		resp2, err := http.DefaultClient.Get(server.URL)
		require.NoError(t, err)
		defer resp2.Body.Close()

		go func() {
			stream.EndAll()
		}()

		data, err = io.ReadAll(resp2.Body)
		require.NoError(t, err)
		assert.Equal(t, "hello,", string(data))
	})
}

func TestChunkedStreamingHandlerClose(t *testing.T) {
	initialData := []byte("hello,")
	handler, stream := ChunkedStreamingHandler(initialData, "text/plain")
	defer stream.Close()

	WithServer(handler, func(server *httptest.Server) {
		resp1, err := http.DefaultClient.Get(server.URL)
		require.NoError(t, err)
		defer resp1.Body.Close()

		go func() {
			stream.Send([]byte("goodbye."))
			stream.Close()
		}()

		data, err := io.ReadAll(resp1.Body)
		require.NoError(t, err)
		assert.Equal(t, "hello,goodbye.", string(data))

		// Should error out on any further requests

		resp2, err := http.DefaultClient.Get(server.URL)
		require.NoError(t, err)
		assert.Equal(t, 500, resp2.StatusCode)
	})
}

func TestChunkedStreamingHandlerWithEnvironmentID(t *testing.T) {
	initialData := []byte("hello")
	handler, stream := ChunkedStreamingHandler(
		initialData,
		"text/plain",
		ChunkedStreamingHandlerOptionEnvironmentID("env-id"),
	)
	defer stream.Close()

	WithServer(handler, func(server *httptest.Server) {
		resp1, err := http.DefaultClient.Get(server.URL)
		require.NoError(t, err)
		defer resp1.Body.Close()

		assert.Equal(t, 200, resp1.StatusCode)
		assert.Equal(t, "text/plain", resp1.Header.Get("Content-Type"))
		assert.Equal(t, "env-id", resp1.Header.Get("X-Ld-Envid"))

		stream.EndAll()

		data, err := io.ReadAll(resp1.Body)
		require.NoError(t, err)
		assert.Equal(t, "hello", string(data))
	})
}
