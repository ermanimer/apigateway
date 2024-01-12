package healthcheckhandler

import "net/http"

type handler struct{}

func New() http.Handler {
	return &handler{}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
