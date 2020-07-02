package httphelpers

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientFromHandler(t *testing.T) {
	handler := HandlerWithStatus(418)
	client := ClientFromHandler(handler)

	resp, err := client.Get("/")

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 418, resp.StatusCode)
}

func TestClientFromHandlerConvertsPanicToError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("sorry")
	})
	client := ClientFromHandler(handler)

	resp, err := client.Get("/")

	expectedError := &url.Error{Op: "Get", URL: "/", Err: errors.New("error from handler: sorry")}
	require.Error(t, err)
	require.Nil(t, resp)
	assert.Equal(t, expectedError, err)
}
