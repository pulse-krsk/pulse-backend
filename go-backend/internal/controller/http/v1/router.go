package v1

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/kurochkinivan/pulskrsk/config"
	"github.com/kurochkinivan/pulskrsk/internal/usecase"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler interface {
	Register(mux *http.ServeMux)
}

// @title			Pulse-krsk API
// @description	pulse kransnoyarsk
// @version		1.0
// @host			localhost:8080
// @BasePath		/api/v1
func NewRouter(cfg *config.Config, a usecase.Auth, u usecase.User) error {
	mux := http.NewServeMux()

	proxyURL := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%s", cfg.JavaService.Host, cfg.JavaService.Port),
	}

	authHandler := NewAuthHandler(a, cfg.BytesLimit, cfg.JWTSignKey)
	authHandler.Register(mux)

	eventHandler := NewProxyHandler(proxyURL)
	eventHandler.Register(mux)

	userHandler := NewUserHandler(u, cfg.BytesLimit, cfg.JWTSignKey)
	userHandler.Register(mux)

	httpSwagger.Handler()
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port), mux)
}
