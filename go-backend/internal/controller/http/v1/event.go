package v1

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kurochkinivan/pulskrsk/config"
)

type eventHandler struct {
	cfg config.JavaService
}

func NewEventHandler(cfg config.JavaService) *eventHandler {
	return &eventHandler{
		cfg: cfg,
	}
}

// TODO: change to httputil.proxy
func (h *eventHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("%s %s/events", http.MethodGet, basePath), errMdw(logMdw(h.getEvents)))
	mux.HandleFunc(fmt.Sprintf("%s %s/events/:id", http.MethodGet, basePath), errMdw(logMdw(h.getEvents)))
}

func (h *eventHandler) getEvents(w http.ResponseWriter, r *http.Request) error {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%s", h.cfg.Host, h.cfg.Port)})
	proxy.ServeHTTP(w, r)

	return nil
}

func (h *eventHandler) getEvent(w http.ResponseWriter, r *http.Request) error {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%s", h.cfg.Host, h.cfg.Port)})
	proxy.ServeHTTP(w, r)

	return nil
}
