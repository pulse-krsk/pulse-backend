package usecase

import "github.com/kurochkinivan/pulskrsk/config"

type UseCases struct {
	Auth
}

type UseCasesDependencies struct {
	UserRepo    UserRepository
	RefreshRepo RefreshSessionsRepository
	*config.Config
}

func NewUseCases(d UseCasesDependencies) *UseCases {
	return &UseCases{
		Auth: NewAuthUseCase(d.UserRepo, d.RefreshRepo, d.Config.Auth),
	}
}
