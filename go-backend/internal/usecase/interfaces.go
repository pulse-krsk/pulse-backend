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
		CreateUser(ctx context.Context, user entity.User) (string, error)
		UserExists(ctx context.Context, oauthID string) (bool, error)
		GetUserByOauthID(ctx context.Context, oauthID string) (entity.User, error)
		GetUserByUUID(ctx context.Context, id string) (entity.User, error)
	}

	RefreshSessionsRepository interface {
		CreateRefreshSession(ctx context.Context, refreshSession entity.RefreshSession) (string, error)
		GetRefreshSession(ctx context.Context, refreshToken string) (entity.RefreshSession, error)
		DeleteRefreshSessionByToken(ctx context.Context, refreshToken string) error
	}

	RoleRepository interface {
		CreateRole(ctx context.Context, role entity.Role) error
		GetRolesByUserID(ctx context.Context, userID string) ([]entity.Role, error)
	}
)

type (
	User interface {
		AddFavoriteEventTypes(ctx context.Context, userID string, eventTypes []string) error
		GetUserWithTypes(ctx context.Context, userID string) (entity.UserWithTypes, error)
	}

	UserFavoriteTypesRepository interface {
		AddUserFavouriteTypes(ctx context.Context, userID string, typeIDs []string) error
		DeleteAllUserFavouriteTypes(ctx context.Context, userID string) error
		GetAllUserFavouriteTypesID(ctx context.Context, userID string) ([]string, error)
	}

	EventTypesRepository interface {
		GetEventTypeByID(ctx context.Context, typeID string) (entity.EventType, error)
		GetEventTypeByType(ctx context.Context, eventType string) (entity.EventType, error)
	}
)
