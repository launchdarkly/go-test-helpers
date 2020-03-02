package httphelpers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithServer(t *testing.T) {
	handler := HandlerWithStatus(200)
	var url string
	WithServer(handler, func(server *httptest.Server) {
		url = server.URL
		resp, err := http.DefaultClient.Get(url)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode)
	})
	_, err := http.DefaultClient.Get(url)
	require.Error(t, err) // server is no longer listening
}
