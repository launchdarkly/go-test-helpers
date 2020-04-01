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
	"gopkg.in/launchdarkly/go-sdk-common.v1/ldvalue"
)

func TestStreamingEndpoint(t *testing.T) {
	initialEvent := NewSSEEvent("", "put", "initial-data")
	eventsCh := make(chan eventsource.Event)
	defer close(eventsCh)
	handler, closer := StreamingServiceHandler(initialEvent, eventsCh)
	defer closer.Close()

	// Note that we have to use an httptest.Server instead of a hacked http.Client because eventsource does not
	// work correctly with httptest.ResponseRecorder.
	httphelpers.WithServer(handler, func(server *httptest.Server) {
		stream, err := eventsource.SubscribeWithURL(server.URL + "/any-url-path")
		require.NoError(t, err)
		defer stream.Close()

		event1 := <-stream.Events
		assert.Equal(t, "put", event1.Event())
		assert.Equal(t, initialEvent.Data(), event1.Data())

		event2 := NewSSEEvent("my-id", "my-event", "my-data")
		eventsCh <- event2

		event3 := <-stream.Events
		assert.Equal(t, event2.Id(), event3.Id())
		assert.Equal(t, event2.Event(), event3.Event())
		assert.Equal(t, event2.Data(), event3.Data())
	})
}

func TestStreamingEndpointClosesStreamWhenHandlerIsClosed(t *testing.T) {
	data1 := NewSSEEvent("", "put", "data1")
	data2 := NewSSEEvent("", "put", "data2")

	// Set up a stream handler that sends data1
	handler1, closer1 := StreamingServiceHandler(data1, nil)
	defer closer1.Close()

	// Set up a stream handler that sends data2
	handler2, closer2 := StreamingServiceHandler(data2, nil)
	defer closer2.Close()

	// Make it so the first request will get handler1, and the second request will get handler2
	sequentialHandler := httphelpers.SequentialHandler(handler1, handler2)

	httphelpers.WithServer(sequentialHandler, func(server *httptest.Server) {
		stream, err := eventsource.SubscribeWithURL(server.URL,
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

func TestServerSideStreamingEndpoint(t *testing.T) {
	data := NewServerSDKData().Flags(FlagOrSegment("flagkey", 1))
	eventsCh := make(chan eventsource.Event)
	defer close(eventsCh)
	handler, closer := ServerSideStreamingServiceHandler(data, eventsCh)
	defer closer.Close()

	httphelpers.WithServer(handler, func(server *httptest.Server) {
		stream, err := eventsource.SubscribeWithURL(server.URL + ServerSideSDKStreamingPath)
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

func TestServerSideStreamingEndpointReturns404ForWrongURL(t *testing.T) {
	data := NewServerSDKData().Flags(FlagOrSegment("flagkey", 1))
	handler, _ := ServerSideStreamingServiceHandler(data, nil)

	httphelpers.WithServer(handler, func(server *httptest.Server) {
		resp, _ := http.DefaultClient.Get(server.URL + "/some/other/path")
		assert.Equal(t, 404, resp.StatusCode)
	})
}

func TestServerSideStreamingEndpointReturns405ForWrongMethod(t *testing.T) {
	eventsCh := make(chan eventsource.Event)
	defer close(eventsCh)
	handler, _ := ServerSideStreamingServiceHandler(nil, eventsCh)

	httphelpers.WithServer(handler, func(server *httptest.Server) {
		resp, _ := http.DefaultClient.Post(server.URL+ServerSideSDKStreamingPath, "text/plain", bytes.NewBufferString("hello"))
		assert.Equal(t, 405, resp.StatusCode)
	})
}

func TestClientSideStreamingEndpoint(t *testing.T) {
	data := NewClientSDKData().Flags(FlagValueData{Key: "flagkey", Version: 1, Value: ldvalue.String("x")})
	eventsCh := make(chan eventsource.Event)
	defer close(eventsCh)
	handler, closer := ClientSideStreamingServiceHandler(data, eventsCh)
	defer closer.Close()

	httphelpers.WithServer(handler, func(server *httptest.Server) {
		stream1, err := eventsource.SubscribeWithURL(server.URL + ClientSideSDKStreamingBasePath + "/envxxx/userxxx")
		require.NoError(t, err)
		defer stream1.Close()

		event1 := <-stream1.Events
		assert.Equal(t, "put", event1.Event())
		assert.Equal(t, data.Data(), event1.Data())

		event2 := FlagValueData{Key: "flagkey", Version: 2, Value: ldvalue.String("y")}
		eventsCh <- event2

		event3 := <-stream1.Events
		assert.Equal(t, event2.Id(), event3.Id())
		assert.Equal(t, event2.Event(), event3.Event())
		assert.Equal(t, event2.Data(), event3.Data())

		stream2, err := eventsource.SubscribeWithURL(server.URL + MobileSDKStreamingBasePath + "/userxxx")
		require.NoError(t, err)
		stream2.Close()

		req3, _ := http.NewRequest("REPORT", server.URL+ClientSideSDKStreamingBasePath+"/envxxx", bytes.NewReader([]byte("{}")))
		stream3, err := eventsource.SubscribeWithRequestAndOptions(req3)
		require.NoError(t, err)
		stream3.Close()

		req4, _ := http.NewRequest("REPORT", server.URL+MobileSDKStreamingBasePath, bytes.NewReader([]byte("{}")))
		stream4, err := eventsource.SubscribeWithRequestAndOptions(req4)
		require.NoError(t, err)
		stream4.Close()
	})
}
