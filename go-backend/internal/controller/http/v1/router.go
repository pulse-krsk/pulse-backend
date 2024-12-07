package v1

import (
	"fmt"
	"net/http"

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

func NewRouter(host, port string, bytesLimit int64, sigingKey string, a usecase.Auth) error {
	mux := http.NewServeMux()

	authHandler := NewAuthHandler(a, bytesLimit, sigingKey)
	authHandler.Register(mux)

	return http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), mux)
}
