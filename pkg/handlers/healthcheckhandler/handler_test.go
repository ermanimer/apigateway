package healthcheckhandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServeHTTP(t *testing.T) {
	expected := http.StatusOK
	handler := New()
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)
	actual := responseRecorder.Code
	require.Equal(t, expected, actual)
}
