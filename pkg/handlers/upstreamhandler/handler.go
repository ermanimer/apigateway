package upstreamhandler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/ermanimer/apigateway/pkg/config"
)

func New(c config.Upstream) http.Handler {
	url, _ := url.Parse(c.URL)
	var handler http.Handler = httputil.NewSingleHostReverseProxy(url)
	if c.StripPrefix {
		prefix := strings.TrimSuffix(c.Pattern, "/")
		handler = http.StripPrefix(prefix, handler)
	}
	return handler
}
