package upstreamhandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ermanimer/apigateway/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	expected := http.StatusOK
	server1 := newTestServer("/server1/endpoint", expected)
	defer server1.Close()
	server2 := newTestServer("/endpoint", expected)
	defer server2.Close()
	tests := []struct {
		name     string
		upstream config.Upstream
		address  string
	}{
		{
			name: "with prefix",
			upstream: config.Upstream{
				Pattern:     "/server1/",
				StripPrefix: false,
				URL:         server1.URL,
			},
			address: "/server1/endpoint",
		},
		{
			name: "strip prefix",
			upstream: config.Upstream{
				Pattern:     "/server2/",
				StripPrefix: true,
				URL:         server2.URL,
			},
			address: "/server2/endpoint",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New(tt.upstream)
			request := httptest.NewRequest(http.MethodGet, tt.address, http.NoBody)
			responseRecorder := httptest.NewRecorder()
			handler.ServeHTTP(responseRecorder, request)
			actual := responseRecorder.Result().StatusCode
			require.Equal(t, expected, actual)
		})
	}
}

type mockHandler struct {
	statusCode int
}

func newMockHandler(statusCode int) http.Handler {
	return &mockHandler{
		statusCode: statusCode,
	}
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(h.statusCode)
}

func newTestServer(path string, statusCode int) *httptest.Server {
	handler := newMockHandler(statusCode)
	mux := http.NewServeMux()
	mux.Handle(path, handler)
	server := httptest.NewServer(mux)
	return server
}
