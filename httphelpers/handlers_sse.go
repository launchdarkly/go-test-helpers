package httphelpers

import (
	"bytes"
	"fmt"
	"net/http"
)

// SSEEvent is a simple representation of a Server-Sent Events message.
type SSEEvent struct {
	// ID is the optional unique ID of the event.
	ID string

	// Event is the message type, if any.
	Event string

	// Data is the message data.
	Data string
}

// Bytes returns the stream data for the event.
func (e SSEEvent) Bytes() []byte {
	var buf bytes.Buffer
	if e.ID != "" {
		buf.WriteString(fmt.Sprintf("id: %s\n", e.ID))
	}
	if e.Event != "" {
		buf.WriteString(fmt.Sprintf("event: %s\n", e.Event))
	}
	buf.WriteString(fmt.Sprintf("data: %s\n\n", e.Data))
	return buf.Bytes()
}

// SSEStreamControl is the interface for manipulating streams created by SSEHandler.
type SSEStreamControl interface {
	// Enqueue is the same as Send, except that if there are currently no open connections to this
	// endpoint, the event is enqueued and will be sent to the next client that connects.
	Enqueue(event SSEEvent)

	// Send sends an SSE event. If there are multiple open connections to this endpoint, the same
	// event is sent to all of them. If there are no open connections, the event is discarded.
	Send(event SSEEvent)

	// EndAll terminates any existing connections to this endpoint, but allows new connections
	// afterward.
	EndAll()

	// Close terminates any existing connections to this endpoint and causes the handler to reject any
	// subsequent connections.
	Close() error
}

type sseStreamControlImpl struct {
	streamControl StreamControl
}

// SSEHandler creates an HTTP handler that streams Server-Sent Events data.
//
// The initialData parameter, if not nil, specifies a starting event that should always be sent to any
// client that has connected to this endpoint. Then, any data provided via the SSEStreamControl interface
// is copied to all connected clients. Connections remain open until either EndAll or Close is called on
// on the SSEStreamControl.
//
// In this example, every request to this endpoint will receive an initial initEvent, and then another
// event will be sent every second with a counter; every 30 seconds, all active stream connections are
// are closed:
//
//     initialEvent := httphelpers.SSEEvent{Data: "hello"}
//     handler, stream := httphelpers.SSEHandler(&initialEvent)
//     (start server with handler)
//     go func() {
//         n := 1
//         counter := time.NewTicker(time.Second)
//         interrupter := time.NewTicker(time.Second * 10)
//         for {
//             select {
//             case <-counter.C:
//                 stream.Send(httphelpers.SSEEvent{Data: fmt.Sprintf("%d\n", n)})
//                 n++
//             case <-interrupter.C:
//                 stream.EndAll()
//             }
//         }
//     }()
func SSEHandler(initialEvent *SSEEvent) (http.Handler, SSEStreamControl) {
	var initialData []byte
	if initialEvent != nil {
		initialData = initialEvent.Bytes()
	}
	handler, streamControl := ChunkedStreamingHandler(initialData, "text/event-stream; charset=utf-8")
	return handler, &sseStreamControlImpl{streamControl}
}

func (s *sseStreamControlImpl) Enqueue(event SSEEvent) {
	s.streamControl.Enqueue(event.Bytes())
}

func (s *sseStreamControlImpl) Send(event SSEEvent) {
	s.streamControl.Send(event.Bytes())
}

func (s *sseStreamControlImpl) EndAll() {
	s.streamControl.EndAll()
}

func (s *sseStreamControlImpl) Close() error {
	return s.streamControl.Close()
}
