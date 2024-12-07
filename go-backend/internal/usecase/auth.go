package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kurochkinivan/pulskrsk/config"

	cuserr "github.com/kurochkinivan/pulskrsk/internal/customErrors"
	"github.com/kurochkinivan/pulskrsk/internal/entity"
	"github.com/kurochkinivan/pulskrsk/internal/usecase/external"
	"github.com/sirupsen/logrus"
)

type AuthUseCase struct {
	userRepo        UserRepository
	refreshRepo     RefreshSessionsRepository
	roleRepo        RoleRepository
	signingKey      string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthUseCase(userRepo UserRepository, refreshRepo RefreshSessionsRepository, roleRepo RoleRepository, auth config.Auth) *AuthUseCase {
	return &AuthUseCase{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		roleRepo:    roleRepo,

		signingKey:      auth.JWTSignKey,
		accessTokenTTL:  auth.AccessTokenTTL,
		refreshTokenTTL: auth.RefreshTokenTTL,
	}
}

func (a *AuthUseCase) LoginUser(ctx context.Context, OauthToken string) (accessToken string, refreshToken string, user entity.User, err error) {
	logrus.WithField("oauth_token", OauthToken).Debug("logging user in")
	const op string = "AuthUseCase.LoginUser"

	user, err = external.ParseOauthToken(OauthToken)
	if err != nil {
		return "", "", entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	exists, err := a.userRepo.UserExists(ctx, user.OauthID)
	if err != nil {
		return "", "", entity.User{}, cuserr.SystemError(err, op, "failed to check if user exists")
	}
	if !exists {
		userID, err := a.userRepo.CreateUser(ctx, user)
		if err != nil {
			return "", "", entity.User{}, cuserr.SystemError(err, op, "failed to create user")
		}

		role := entity.Role{UserID: userID, RoleName: "user"}
		fmt.Println(role)
		err = a.roleRepo.CreateRole(ctx, role)
		if err != nil {
			return "", "", entity.User{}, cuserr.SystemError(err, op, "failed to create user role")
		}
	}

	user, err = a.userRepo.GetUserByOauthID(ctx, user.OauthID)
	if err != nil {
		return "", "", entity.User{}, cuserr.SystemError(err, op, "failed to get user")
	}

	roles, err := a.roleRepo.GetRolesByUserID(ctx, user.ID)
	if err != nil {
		return "", "", entity.User{}, cuserr.SystemError(err, op, "failed to get user role")
	}

	accessToken, refreshToken, err = a.GenerateTokenPair(ctx, user.ID, roles)
	if err != nil {
		return "", "", entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return accessToken, refreshToken, user, nil
}

func (a *AuthUseCase) LogoutUser(ctx context.Context, refreshToken string) error {
	logrus.WithField("refresh_token", refreshToken).Debug("logging user out")
	const op string = "authUseCase.LogoutUser"

	err := a.refreshRepo.DeleteRefreshSessionByToken(ctx, refreshToken)
	if err != nil {
		return cuserr.SystemError(err, op, "failed to delete refresh session")
	}

	return nil
}

func (a *AuthUseCase) RefreshTokens(ctx context.Context, refreshTkn string) (accessToken string, refreshToken string, err error) {
	logrus.WithField("refresh_token", refreshToken).Debug("refreshing tokens")
	const op string = "AuthUseCase.RefreshTokens"

	refreshSession, err := a.refreshRepo.GetRefreshSession(ctx, refreshTkn)
	if err != nil {
		return "", "", cuserr.SystemError(err, op, "failed to get refresh session")
	}
	if refreshSession == (entity.RefreshSession{}) {
		return "", "", cuserr.ErrInvalidRefreshToken
	}

	err = a.refreshRepo.DeleteRefreshSessionByToken(ctx, refreshTkn)
	if err != nil {
		return "", "", cuserr.SystemError(err, op, "failed to delete refresh session by refresh token")
	}

	if refreshSession.Expiration.Before(time.Now()) {
		return "", "", cuserr.ErrTokenExired
	}

	user, err := a.userRepo.GetUserByUUID(ctx, refreshSession.UserID)
	if err != nil {
		return "", "", cuserr.SystemError(err, op, "failed to get user")
	}

	roles, err := a.roleRepo.GetRolesByUserID(ctx, user.ID)
	if err != nil {
		return "", "", cuserr.SystemError(err, op, "failed to get user role")
	}

	accessToken, refreshToken, err = a.GenerateTokenPair(ctx, user.ID, roles)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return accessToken, refreshToken, nil
}

func (a *AuthUseCase) GenerateTokenPair(ctx context.Context, userID string, roles []entity.Role) (accessToken, refreshToken string, err error) {
	logrus.WithField("userID", userID).Debug("generating token pair")
	const op string = "AuthUseCase.GenerateTokenPair"

	accesstoken, err := a.GenerateAccessToken(userID, roles)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	refreshToken, err = a.GenerateRefreshToken(ctx, userID)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return accesstoken, refreshToken, nil
}

func (a *AuthUseCase) GenerateAccessToken(userID string, roles []entity.Role) (string, error) {
	logrus.WithField("id", userID).Debug("generating access token")
	const op string = "AuthUseCase.GenerateAccessToken"

	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.RoleName) // Используем поле Name
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ueid":  userID,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(a.accessTokenTTL).Unix(),
		"roles": roleNames,
	})

	signedToken, err := token.SignedString([]byte(a.signingKey))
	if err != nil {
		logrus.WithError(err).Error("failed to sign token")

		return "", cuserr.SystemError(err, op, "failed to sign token")
	}

	return signedToken, nil
}

func (a *AuthUseCase) GenerateRefreshToken(ctx context.Context, userID string) (string, error) {
	logrus.WithField("userID", userID).Debug("generating refresh token")
	const op string = "AuthUseCase.GenerateRefreshToken"

	refreshSession := entity.RefreshSession{
		UserID:     userID,
		IssuedAt:   time.Now(),
		Expiration: time.Now().Add(a.refreshTokenTTL),
	}
	refreshToken, err := a.refreshRepo.CreateRefreshSession(ctx, refreshSession)
	if err != nil {
		return "", cuserr.SystemError(err, op, "failed to create refresh session")
	}

	return refreshToken, nil
}
