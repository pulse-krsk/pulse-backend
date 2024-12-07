package usecase

import "github.com/kurochkinivan/pulskrsk/config"

type UseCases struct {
	Auth
	User
}

type UseCasesDependencies struct {
	UserRepo             UserRepository
	RefreshRepo          RefreshSessionsRepository
	RoleRepo             RoleRepository
	EventTypesRepo       EventTypesRepository
	UserFavoriteTypeRepo UserFavoriteTypesRepository
	*config.Config
}

func NewUseCases(d UseCasesDependencies) *UseCases {
	return &UseCases{
		Auth: NewAuthUseCase(d.UserRepo, d.RefreshRepo, d.RoleRepo, d.Config.Auth),
		User: NewUserUseCase(d.UserFavoriteTypeRepo, d.EventTypesRepo, d.UserRepo),
	}
}
