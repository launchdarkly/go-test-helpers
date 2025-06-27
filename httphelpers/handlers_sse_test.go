package httphelpers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSEHandler(t *testing.T) {
	initialEvent := SSEEvent{"id1", "event1", "data1", 0}
	handler, stream := SSEHandler(&initialEvent)
	defer stream.Close()

	stream.Enqueue(SSEEvent{"", "event2", "data2", 0})
	stream.EnqueueComment("comment1")
	stream.Send(SSEEvent{"", "", "this isn't sent because there are no connections", 0})

	WithServer(handler, func(server *httptest.Server) {
		resp1, err := http.DefaultClient.Get(server.URL)
		require.NoError(t, err)
		defer resp1.Body.Close()

		assert.Equal(t, 200, resp1.StatusCode)
		assert.Equal(t, "text/event-stream; charset=utf-8", resp1.Header.Get("Content-Type"))

		stream.SendComment("comment2")
		stream.Enqueue(SSEEvent{"", "event3", "data3", 500})
		stream.EndAll()

		data, err := io.ReadAll(resp1.Body)

		assert.NoError(t, err)
		assert.Equal(t, `id: id1
event: event1
data: data1

event: event2
data: data2

:comment1
:comment2
event: event3
retry: 500
data: data3

`, string(data))
	})
}

func TestSSEHandlerWithEnvironmentID(t *testing.T) {
	initialEvent := SSEEvent{"id1", "event1", "data1", 0}
	handler, stream := SSEHandlerWithEnvironmentID(&initialEvent, "env-id")
	defer stream.Close()

	WithServer(handler, func(server *httptest.Server) {
		resp1, err := http.DefaultClient.Get(server.URL)
		require.NoError(t, err)
		defer resp1.Body.Close()

		assert.Equal(t, 200, resp1.StatusCode)
		assert.Equal(t, "text/event-stream; charset=utf-8", resp1.Header.Get("Content-Type"))
		assert.Equal(t, "env-id", resp1.Header.Get("X-Ld-Envid"))

		stream.EndAll()

		data, err := io.ReadAll(resp1.Body)
		assert.NoError(t, err)
		assert.Equal(t, `id: id1
event: event1
data: data1

`, string(data))
	})
}
