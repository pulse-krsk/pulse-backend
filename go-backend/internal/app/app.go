package app

import (
	"context"
	"fmt"

	"github.com/kurochkinivan/pulskrsk/config"
	v1 "github.com/kurochkinivan/pulskrsk/internal/controller/http/v1"
	"github.com/kurochkinivan/pulskrsk/internal/usecase"
	"github.com/kurochkinivan/pulskrsk/internal/usecase/repository/postgresql"
	psql "github.com/kurochkinivan/pulskrsk/pkg/postgresql"
	"github.com/sirupsen/logrus"
)

func Run(cfg *config.Config) error {

	logrus.Info("connecting to database client...")
	cfgp := cfg.PostgreSQL
	pgConfig := psql.NewPgConfig(cfgp.Username, cfgp.Password, cfgp.Host, cfgp.Port, cfgp.Database)
	client, err := psql.NewClient(context.Background(), 5, pgConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
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
	return v1.NewRouter(cfg, authUseCase, userUseCase)
}
