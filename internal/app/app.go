package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/kurochkinivan/pulskrsk/config"
	v1 "github.com/kurochkinivan/pulskrsk/internal/controller/http/v1"
	"github.com/kurochkinivan/pulskrsk/internal/usecase"
	"github.com/kurochkinivan/pulskrsk/internal/usecase/repository/postgresql"
	psql "github.com/kurochkinivan/pulskrsk/pkg/postgresql"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

type App struct {
	cfg        *config.Config
	router     *http.ServeMux
	httpServer *http.Server
}

func NewApp(cfg *config.Config) (*App, error) {
	logrus.Info("connecting to database client...")
	cfgp := cfg.PostgreSQL
	pgConfig := psql.NewPgConfig(cfgp.Username, cfgp.Password, cfgp.Host, cfgp.Port, cfgp.Database)
	client, err := psql.NewClient(context.Background(), 5, pgConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logrus.Info("creating repositories and usecases...")
	authRepo := postgresql.NewUserRepository(client)
	refreshRepo := postgresql.NewRefreshSessionsRepository(client)
	roleRepo := postgresql.NewRoleRepository(client)
	eventTypeRepo := postgresql.NewEventTypeRepository(client)
	usersFavoriteTypesRepo := postgresql.NewUsersFavoriteTypesRepository(client)
	dependencies := usecase.UseCasesDependencies{
		UserRepo:             authRepo,
		RefreshRepo:          refreshRepo,
		RoleRepo:             roleRepo,
		UserFavoriteTypeRepo: usersFavoriteTypesRepo,
		EventTypesRepo:       eventTypeRepo,
		Config:               cfg,
	}
	authUseCase := usecase.NewUseCases(dependencies)
	userUseCase := usecase.NewUseCases(dependencies)

	logrus.Info("starting server...")
	router := v1.NewRouter(cfg, authUseCase, userUseCase)

	app := &App{
		cfg:    cfg,
		router: router,
	}

	return app, nil
}

func (a *App) StartHTTP() error {
	logrus.Info("creating listener")
	logrus.WithFields(logrus.Fields{"IP": a.cfg.HTTP.Host, "Port": a.cfg.HTTP.Port}).Debug("http listener credentials")
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.HTTP.Host, a.cfg.HTTP.Port))
	if err != nil {
		logrus.WithError(err).Fatal("failed to create listener")
	}

	logrus.WithFields(map[string]interface{}{
		"AllowedMethods":     a.cfg.HTTP.CORS.AllowedMethods,
		"AllowedOrigins":     a.cfg.HTTP.CORS.AllowedOrigins,
		"AllowCredentials":   a.cfg.HTTP.CORS.AllowCredentials,
		"AllowedHeaders":     a.cfg.HTTP.CORS.AllowedHeaders,
		"OptionsPassthrough": a.cfg.HTTP.CORS.OptionsPassthrough,
		"ExposedHeaders":     a.cfg.HTTP.CORS.ExposedHeaders,
		"Debug":              a.cfg.HTTP.CORS.Debug,
	})
	c := cors.New(cors.Options{
		AllowedMethods:     a.cfg.HTTP.CORS.AllowedMethods,
		AllowedOrigins:     a.cfg.HTTP.CORS.AllowedOrigins,
		AllowCredentials:   a.cfg.HTTP.CORS.AllowCredentials,
		AllowedHeaders:     a.cfg.HTTP.CORS.AllowedHeaders,
		OptionsPassthrough: a.cfg.HTTP.CORS.OptionsPassthrough,
		ExposedHeaders:     a.cfg.HTTP.CORS.ExposedHeaders,
		Debug:              a.cfg.HTTP.CORS.Debug,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: a.cfg.HTTP.WriteTimeout,
		ReadTimeout:  a.cfg.HTTP.ReadTimeout,
	}

	logrus.Info("application completely initialized and started")

	if err = a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logrus.Warn("server shutdown")
		default:
			logrus.Fatal(err)
		}
	}

	err = a.httpServer.Shutdown(context.Background())
	if err != nil {
		logrus.Fatal(err)
	}

	return err
}
