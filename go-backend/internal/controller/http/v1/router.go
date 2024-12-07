package v1

import (
	"fmt"
	"net/http"

	"github.com/kurochkinivan/pulskrsk/config"
	"github.com/kurochkinivan/pulskrsk/internal/usecase"
)

type Handler interface {
	Register(mux *http.ServeMux)
}

type Handlers struct {
	Auth authHandler
}

func NewHandlers(auth authHandler) *Handlers {
	return &Handlers{
		Auth: auth,
	}
}

func NewRouter(cfg *config.Config, a usecase.Auth) error {
	mux := http.NewServeMux()

	authHandler := NewAuthHandler(a, cfg.BytesLimit, cfg.JWTSignKey)
	authHandler.Register(mux)

	eventHandler := NewEventHandler(cfg.JavaService)
	eventHandler.Register(mux)

	return http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port), mux)
}
