package usecase

import (
	"context"

	cuserr "github.com/kurochkinivan/pulskrsk/internal/customErrors"
	"github.com/kurochkinivan/pulskrsk/internal/entity"
	"github.com/sirupsen/logrus"
)

type userUseCase struct {
	userFavoriteTypesRepo UserFavoriteTypesRepository
	eventTypesRepo        EventTypesRepository
	userRepo              UserRepository
}

func NewUserUseCase(userFavoriteTypesRepo UserFavoriteTypesRepository, eventTypesRepo EventTypesRepository, userRepo UserRepository) *userUseCase {
	return &userUseCase{
		userFavoriteTypesRepo: userFavoriteTypesRepo,
		eventTypesRepo:        eventTypesRepo,
		userRepo:              userRepo,
	}
}

func (u *userUseCase) AddFavoriteEventTypes(ctx context.Context, userID string, eventTypes []string) error {
	logrus.WithField("user_id", userID).Debug("adding favorite event types")
	const op string = "userUseCase.AddFavoriteEventTypes"

	err := u.userFavoriteTypesRepo.DeleteAllUserFavouriteTypes(ctx, userID)
	if err != nil {
		return cuserr.SystemError(err, op, "failed to delete all user favorite event types")
	}

	eventTypesIDs := make([]string, 0, len(eventTypes))
	for _, eventType := range eventTypes {
		eventType, err := u.eventTypesRepo.GetEventTypeByType(ctx, eventType)
		if err != nil {
			return cuserr.SystemError(err, op, "failed to get event by type")
		}
		eventTypesIDs = append(eventTypesIDs, eventType.ID)
	}

	err = u.userFavoriteTypesRepo.AddUserFavouriteTypes(ctx, userID, eventTypesIDs)
	if err != nil {
		return cuserr.SystemError(err, op, "failed to add user favorite event types")
	}

	return nil
}

func (u *userUseCase) GetUserWithTypes(ctx context.Context, userID string) (entity.UserWithTypes, error) {
	logrus.WithField("user_id", userID).Debug("getting user")
	const op string = "userUseCase.GetUserWithTypes"

	user, err := u.userRepo.GetUserByUUID(ctx, userID)
	if err != nil {
		return entity.UserWithTypes{}, cuserr.SystemError(err, op, "failed to get user")
	}

	typeIDs, err := u.userFavoriteTypesRepo.GetAllUserFavouriteTypesID(ctx, user.ID)
	if err != nil {
		return entity.UserWithTypes{}, cuserr.SystemError(err, op, "failed to get user favorite event types")
	}

	types := make([]entity.EventType, 0, len(typeIDs))
	for _, typeID := range typeIDs {
		type_, err := u.eventTypesRepo.GetEventTypeByID(ctx, typeID)
		if err != nil {
			return entity.UserWithTypes{}, cuserr.SystemError(err, op, "failed to get event by id")
		}
		types = append(types, type_)
	}

	var userWithTypes entity.UserWithTypes
	userWithTypes.User = user
	userWithTypes.Types = types

	return userWithTypes, nil
}
