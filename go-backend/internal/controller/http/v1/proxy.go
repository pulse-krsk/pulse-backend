package v1

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type proxyHandler struct {
	proxyURL *url.URL
}

func NewProxyHandler(proxyURL *url.URL) *proxyHandler {
	return &proxyHandler{
		proxyURL: proxyURL,
	}
}

// TODO: change to httputil.proxy
func (h *proxyHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("%s/events", basePath), errMdw(logMdw(h.proxyEvents)))
	mux.HandleFunc(fmt.Sprintf("%s/events/", basePath), errMdw(logMdw(h.proxyEvents)))
	mux.HandleFunc(fmt.Sprintf("%s/event-types", basePath), errMdw(logMdw(h.proxyEvents)))
	mux.HandleFunc(fmt.Sprintf("%s/users/", basePath), errMdw(logMdw(h.proxyEvents)))
}

func (h *proxyHandler) proxyEvents(w http.ResponseWriter, r *http.Request) error {
	proxy := httputil.NewSingleHostReverseProxy(h.proxyURL)
	proxy.ServeHTTP(w, r)

	return nil
}
