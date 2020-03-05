package httphelpers

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelfSignedServer(t *testing.T) {
	handler := HandlerWithStatus(200)
	WithSelfSignedServer(handler, func(server *httptest.Server, certData []byte, certs *x509.CertPool) {
		client := *http.DefaultClient
		transport := &http.Transport{}
		transport.TLSClientConfig = &tls.Config{RootCAs: certs}
		client.Transport = transport
		resp, err := client.Get(server.URL)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode)
	})
}
