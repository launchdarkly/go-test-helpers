package httphelpers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSEHandler(t *testing.T) {
	initialEvent := SSEEvent{"id1", "event1", "data1"}
	handler, stream := SSEHandler(&initialEvent)
	defer stream.Close()

	stream.Enqueue(SSEEvent{"", "event2", "data2"})
	stream.Send(SSEEvent{"", "", "this isn't sent becauset here are no connections"})

	WithServer(handler, func(server *httptest.Server) {
		resp1, err := http.DefaultClient.Get(server.URL)
		require.NoError(t, err)
		defer resp1.Body.Close()

		assert.Equal(t, 200, resp1.StatusCode)
		assert.Equal(t, "text/event-stream; charset=utf-8", resp1.Header.Get("Content-Type"))

		stream.Enqueue(SSEEvent{"", "event3", "data3"})
		stream.EndAll()

		data, err := ioutil.ReadAll(resp1.Body)

		assert.NoError(t, err)
		assert.Equal(t, `id: id1
event: event1
data: data1

event: event2
data: data2

event: event3
data: data3

`, string(data))
	})
}
