package httphelpers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelegatingHandler(t *testing.T) {
	h1 := HandlerWithStatus(200)
	h2 := HandlerWithStatus(304)
	dh := DelegatingHandler{h1}

	rr1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/", nil)
	dh.ServeHTTP(rr1, req1)
	assert.Equal(t, 200, rr1.Code)

	dh.Handler = h2

	rr2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/", nil)
	dh.ServeHTTP(rr2, req2)
	assert.Equal(t, 304, rr2.Code)
}

func TestHandlerForMethod(t *testing.T) {
	h1 := HandlerWithStatus(200)
	h2 := HandlerWithStatus(202)
	hm := HandlerForMethod("GET", h1, HandlerForMethod("POST", h2, nil))

	rr1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/", nil)
	hm.ServeHTTP(rr1, req1)
	assert.Equal(t, 200, rr1.Code)

	rr2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/", nil)
	hm.ServeHTTP(rr2, req2)
	assert.Equal(t, 202, rr2.Code)

	rr3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("PATCH", "/", nil)
	hm.ServeHTTP(rr3, req3)
	assert.Equal(t, 405, rr3.Code)
}

func TestHandlerForPath(t *testing.T) {
	h1 := HandlerWithStatus(200)
	h2 := HandlerWithStatus(304)
	hp := HandlerForPath("/path1", h1, HandlerForPath("/path2", h2, nil))

	rr1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/path1", nil)
	hp.ServeHTTP(rr1, req1)
	assert.Equal(t, 200, rr1.Code)

	rr2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/path2", nil)
	hp.ServeHTTP(rr2, req2)
	assert.Equal(t, 304, rr2.Code)

	rr3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/path3", nil)
	hp.ServeHTTP(rr3, req3)
	assert.Equal(t, 404, rr3.Code)
}

func TestHandlerForPathRegex(t *testing.T) {
	h1 := HandlerWithStatus(200)
	h2 := HandlerWithStatus(304)
	hp := HandlerForPathRegex("^/path[12]$", h1, h2)

	rr1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/path1", nil)
	hp.ServeHTTP(rr1, req1)
	assert.Equal(t, 200, rr1.Code)

	rr2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/path2", nil)
	hp.ServeHTTP(rr2, req2)
	assert.Equal(t, 200, rr2.Code)

	rr3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/path3", nil)
	hp.ServeHTTP(rr3, req3)
	assert.Equal(t, 304, rr3.Code)
}

func TestHandlerWithJSONResponse(t *testing.T) {
	jsonObject := map[string]string{"things": "stuff"}
	headers := make(http.Header)
	headers.Set("X-My-Header", "hello")
	h := HandlerWithJSONResponse(jsonObject, headers)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(t, "hello", rr.Header().Get("X-My-Header"))
	assert.Equal(t, []byte(`{"things":"stuff"}`), rr.Body.Bytes())
}

func TestHandlerWithResponse(t *testing.T) {
	headers1 := make(http.Header)
	headers1.Set("X-My-Header", "hello")
	h1 := HandlerWithResponse(200, headers1, nil)

	rr1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/", nil)
	h1.ServeHTTP(rr1, req1)
	assert.Equal(t, 200, rr1.Code)
	assert.Equal(t, "hello", rr1.Header().Get("X-My-Header"))
	assert.Nil(t, rr1.Body.Bytes())

	headers2 := make(http.Header)
	headers2.Set("Content-Type", "text/plain")
	h2 := HandlerWithResponse(200, headers2, []byte("hello"))

	rr2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/", nil)
	h2.ServeHTTP(rr2, req2)
	assert.Equal(t, 200, rr2.Code)
	assert.Equal(t, "text/plain", rr2.Header().Get("Content-Type"))
	assert.Equal(t, []byte("hello"), rr2.Body.Bytes())
}

func TestHandlerWithStatus(t *testing.T) {
	h := HandlerWithStatus(418) // I'm a teapot

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(rr, req)

	assert.Equal(t, 418, rr.Code)
}

func TestRecordingHandler(t *testing.T) {
	h := HandlerWithStatus(418)
	rh, requestsCh := RecordingHandler(h)

	req1, _ := http.NewRequest("GET", "/1", nil)
	rr1 := httptest.NewRecorder()
	rh.ServeHTTP(rr1, req1)

	postData := []byte("hello")
	req2, _ := http.NewRequest("GET", "/2", bytes.NewBuffer(postData))
	rr2 := httptest.NewRecorder()
	rh.ServeHTTP(rr2, req2)

	assert.Equal(t, 418, rr1.Code)
	assert.Equal(t, 418, rr2.Code)

	assert.Equal(t, 2, len(requestsCh))
	ri1 := <-requestsCh
	assert.Equal(t, req1.URL.Path, ri1.Request.URL.Path)
	assert.Nil(t, ri1.Body)
	ri2 := <-requestsCh
	assert.Equal(t, req2.URL.Path, ri2.Request.URL.Path)
	assert.Equal(t, postData, ri2.Body)
}

func TestSequentialHandler(t *testing.T) {
	h1 := HandlerWithStatus(500)
	h2 := HandlerWithStatus(400)
	sh := SequentialHandler(h1, h2)

	req1, _ := http.NewRequest("GET", "/1", nil)
	rr1 := httptest.NewRecorder()
	sh.ServeHTTP(rr1, req1)

	req2, _ := http.NewRequest("GET", "/2", nil)
	rr2 := httptest.NewRecorder()
	sh.ServeHTTP(rr2, req2)

	req3, _ := http.NewRequest("GET", "/3", nil)
	rr3 := httptest.NewRecorder()
	sh.ServeHTTP(rr3, req3)

	assert.Equal(t, 500, rr1.Code)
	assert.Equal(t, 400, rr2.Code)
	assert.Equal(t, 400, rr3.Code)
}

func TestBrokenConnectionHandler(t *testing.T) {
	h := BrokenConnectionHandler()

	t.Run("with instrumented client", func(t *testing.T) {
		client := ClientFromHandler(h)
		_, err := client.Get("/")
		assert.Error(t, err)
	})

	t.Run("with server", func(t *testing.T) {
		WithServer(h, func(server *httptest.Server) {
			_, err := http.DefaultClient.Get(server.URL)
			assert.Error(t, err)
		})
	})
}
