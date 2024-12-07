package usecase

import "github.com/kurochkinivan/pulskrsk/config"

type UseCases struct {
	Auth
}

type UseCasesDependencies struct {
	UserRepo    UserRepository
	RefreshRepo RefreshSessionsRepository
	RoleRepo    RoleRepository
	*config.Config
}

func NewUseCases(d UseCasesDependencies) *UseCases {
	return &UseCases{
		Auth: NewAuthUseCase(d.UserRepo, d.RefreshRepo, d.RoleRepo, d.Config.Auth),
	}
}
