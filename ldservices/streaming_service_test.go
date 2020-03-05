package ldservices

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/launchdarkly/eventsource"
	"github.com/launchdarkly/go-test-helpers/httphelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerSideStreamingEndpoint(t *testing.T) {
	data := NewServerSDKData().Flags(FlagOrSegment("flagkey", 1))
	eventsCh := make(chan eventsource.Event)
	defer close(eventsCh)
	handler, closer := ServerSideStreamingServiceHandler(data, eventsCh)
	defer closer.Close()

	// Note that we have to use an httptest.Server instead of a hacked http.Client because eventsource does not
	// work correctly with httptest.ResponseRecorder.
	httphelpers.WithServer(handler, func(server *httptest.Server) {
		stream, err := eventsource.SubscribeWithURL(server.URL + serverSideSDKStreamingPath)
		require.NoError(t, err)
		defer stream.Close()

		event1 := <-stream.Events
		assert.Equal(t, "put", event1.Event())
		assert.Equal(t, data.Data(), event1.Data())

		event2 := NewSSEEvent("my-id", "my-event", "my-data")
		eventsCh <- event2

		event3 := <-stream.Events
		assert.Equal(t, event2.Id(), event3.Id())
		assert.Equal(t, event2.Event(), event3.Event())
		assert.Equal(t, event2.Data(), event3.Data())
	})
}

func TestServerSideStreamingEndpointClosesStreamWhenHandlerIsClosed(t *testing.T) {
	data1 := NewServerSDKData().Flags(FlagOrSegment("flagkey1", 1))
	data2 := NewServerSDKData().Flags(FlagOrSegment("flagkey2", 2))

	// Set up a stream handler that sends data1
	handler1, closer1 := ServerSideStreamingServiceHandler(data1, nil)
	defer closer1.Close()

	// Set up a stream handler that sends data2
	handler2, closer2 := ServerSideStreamingServiceHandler(data2, nil)
	defer closer2.Close()

	// Make it so the first request will get handler1, and the second request will get handler2
	sequentialHandler := httphelpers.SequentialHandler(handler1, handler2)

	httphelpers.WithServer(sequentialHandler, func(server *httptest.Server) {
		stream, err := eventsource.SubscribeWithURL(server.URL+serverSideSDKStreamingPath,
			eventsource.StreamOptionInitialRetry(time.Millisecond))
		require.NoError(t, err)
		defer stream.Close()

		// Make sure consume stream errors so the error channel won't block on disconnect
		go func() {
			for range stream.Errors {
			}
		}()

		// Wait for the event from handler1
		event1 := <-stream.Events
		assert.Equal(t, "put", event1.Event())
		assert.Equal(t, data1.Data(), event1.Data())

		// Close handler1, so it closes the stream, so eventsource reconnects and gets handler2
		err = closer1.Close()
		require.NoError(t, err)

		event2 := <-stream.Events
		assert.Equal(t, "put", event2.Event())
		assert.Equal(t, data2.Data(), event2.Data())
	})
}

func TestServerSideStreamingEndpointReturns404ForWrongURL(t *testing.T) {
	data := NewServerSDKData().Flags(FlagOrSegment("flagkey", 1))
	handler, _ := ServerSideStreamingServiceHandler(data, nil)

	httphelpers.WithServer(handler, func(server *httptest.Server) {
		resp, _ := http.DefaultClient.Get(server.URL + "/some/other/path")
		assert.Equal(t, 404, resp.StatusCode)
	})
}

func TestServerSideStreamingEndpointReturns405ForWrongMethod(t *testing.T) {
	data := NewServerSDKData().Flags(FlagOrSegment("flagkey", 1))
	eventsCh := make(chan eventsource.Event)
	defer close(eventsCh)
	handler, _ := ServerSideStreamingServiceHandler(data, eventsCh)

	httphelpers.WithServer(handler, func(server *httptest.Server) {
		resp, _ := http.DefaultClient.Post(server.URL+serverSideSDKStreamingPath, "text/plain", bytes.NewBufferString("hello"))
		assert.Equal(t, 405, resp.StatusCode)
	})
}
