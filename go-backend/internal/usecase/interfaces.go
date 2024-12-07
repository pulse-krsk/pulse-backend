package usecase

import (
	"context"

	"github.com/kurochkinivan/pulskrsk/internal/entity"
)

// TODO: think of generate mocks
//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	Auth interface {
		LoginUser(ctx context.Context, OauthToken string) (accessToken string, refreshToken string, user entity.User, err error)
		RefreshTokens(ctx context.Context, refreshTkn string) (accessToken string, refreshToken string, err error)
		LogoutUser(ctx context.Context, refreshToken string) error
	}

	UserRepository interface {
		CreateUser(ctx context.Context, user entity.User) error
		UserExists(ctx context.Context, oauthID string) (bool, error)
		GetUserByOauthID(ctx context.Context, oauthID string) (entity.User, error)
		GetUserByUUID(ctx context.Context, id string) (entity.User, error)
	}

	RefreshSessionsRepository interface {
		CreateRefreshSession(ctx context.Context, refreshSession entity.RefreshSession) (string, error)
		GetRefreshSession(ctx context.Context, refreshToken string) (entity.RefreshSession, error)
		DeleteRefreshSessionByToken(ctx context.Context, refreshToken string) error
		DeleteRefreshSessionsByUserID(ctx context.Context, userID string) error
	}
)
