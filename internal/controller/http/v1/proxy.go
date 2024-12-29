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

func (h *proxyHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc(baseEventPath, errMdw(logMdw(h.proxyEvents)))
	mux.HandleFunc(fmt.Sprintf("%s/", baseEventPath), errMdw(logMdw(h.proxyEvents)))
	mux.HandleFunc(fmt.Sprintf("%s/event-types", basePath), errMdw(logMdw(h.proxyEvents)))
	mux.HandleFunc(fmt.Sprintf("%s/", baseUsersPath), errMdw(logMdw(h.proxyEvents)))
}

func (h *proxyHandler) proxyEvents(w http.ResponseWriter, r *http.Request) error {
	proxy := httputil.NewSingleHostReverseProxy(h.proxyURL)
	proxy.ServeHTTP(w, r)

	return nil
}
