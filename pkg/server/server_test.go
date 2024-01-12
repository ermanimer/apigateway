package server

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ermanimer/apigateway/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	address := ":8080"
	readTimeout := 5 * time.Second
	writeTimeout := 10 * time.Second
	idleTimeout := 120 * time.Second
	maxHeaderBytes := 1048576
	shutdownTimeout := 10 * time.Second
	config := config.Server{
		Address:         address,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		IdleTimeout:     idleTimeout,
		MaxHeaderBytes:  maxHeaderBytes,
		ShutdownTimeout: shutdownTimeout,
	}
	expected := &Server{
		server: &http.Server{
			Addr:           address,
			ReadTimeout:    readTimeout,
			WriteTimeout:   writeTimeout,
			IdleTimeout:    idleTimeout,
			MaxHeaderBytes: maxHeaderBytes,
			Handler:        http.NewServeMux(),
		},
		shutdownTimeout: shutdownTimeout,
	}
	actual := New(config)
	require.Equal(t, expected.server.Addr, actual.server.Addr)
	require.Equal(t, expected.server.ReadTimeout, actual.server.ReadTimeout)
	require.Equal(t, expected.server.WriteTimeout, actual.server.WriteTimeout)
	require.Equal(t, expected.server.IdleTimeout, actual.server.IdleTimeout)
	require.Equal(t, expected.server.MaxHeaderBytes, actual.server.MaxHeaderBytes)
	require.Equal(t, expected.server.Handler, actual.server.Handler)
	require.Equal(t, expected.shutdownTimeout, actual.shutdownTimeout)
}

func TestRegisterHandler(t *testing.T) {
	// create and register a mock handler
	pattern := "/"
	expected := http.StatusOK
	handler, _ := newMockHandler(expected)
	server := New(config.Server{})
	server.RegisterHandler(pattern, handler)

	// create and handle a request with the registered handler
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	responseRecorder := httptest.NewRecorder()
	server.server.Handler.ServeHTTP(responseRecorder, request)
	actual := responseRecorder.Code
	require.Equal(t, expected, actual)
}

func TestStartAndStop(t *testing.T) {
	// get a test address
	address, err := getTestAddress()
	require.NoError(t, err)

	// create a server
	shutdownTimeout := 1 * time.Second
	server := New(config.Server{
		Address:         address,
		ShutdownTimeout: shutdownTimeout,
	})

	// create and register a mock handler
	expected := http.StatusOK
	path := "/"
	handler, done := newMockHandler(expected)
	server.RegisterHandler(path, handler)

	// start the server
	go func() {
		err := server.Start()
		require.ErrorIs(t, err, http.ErrServerClosed)
	}()

	// wait for the server to start, send a request and check the response
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			request, err := http.NewRequest(http.MethodGet, "http://"+address+path, http.NoBody)
			require.NoError(t, err)
			var response *http.Response
			response, err = http.DefaultClient.Do(request)
			if err != nil {
				require.ErrorIs(t, err, net.ErrClosed)
				continue
			}
			require.Equal(t, expected, response.StatusCode)
		}
	}()

	// wait for the request to be handled
	<-done

	// stop the server and check the shutdown duration
	err = server.Shutdown()
	require.NoError(t, err)
}

// mock handler
type mockHandler struct {
	statusCode int
	done       chan struct{}
}

func newMockHandler(statusCode int) (http.Handler, chan struct{}) {
	done := make(chan struct{})
	handler := &mockHandler{
		statusCode: statusCode,
		done:       done,
	}
	return handler, done
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(h.statusCode)
	close(h.done)
}

func getTestAddress() (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		return "", err
	}
	defer listener.Close()
	return listener.Addr().String(), nil
}
