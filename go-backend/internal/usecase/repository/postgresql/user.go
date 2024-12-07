package postgresql

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kurochkinivan/pulskrsk/internal/entity"
	psql "github.com/kurochkinivan/pulskrsk/pkg/postgresql"
	"github.com/sirupsen/logrus"
)

type userRepository struct {
	client psql.PosgreSQLClient
	qb     sq.StatementBuilderType
}

func NewUserRepository(client *pgxpool.Pool) *userRepository {
	return &userRepository{
		client: client,
		qb:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user entity.User) error {
	logrus.WithField("email", user.Email).Trace("creating user")
	const op string = "userRepository.CreateUser"

	uuid, _ := uuid.NewUUID()
	sql, args, err := r.qb.
		Insert(TableUsers).
		Columns(
			"id",
			"first_name",
			"last_name",
			"email",
			"oauth_id",
		).
		Values(
			uuid.String(),
			user.FirstName,
			user.LastName,
			user.Email,
			user.OauthID,
		).
		ToSql()
	if err != nil {
		return psql.ErrCreateQuery(op, err)
	}

	commtag, err := r.client.Exec(ctx, sql, args...)
	if err != nil {
		return psql.ErrExec(op, err)
	}

	if commtag.RowsAffected() == 0 {
		return psql.NoRowsAffected
	}

	return nil
}

func (r *userRepository) UserExists(ctx context.Context, oauthID string) (bool, error) {
	logrus.WithField("oauth_id", oauthID).Trace("checking if user exists")
	const op string = "userRepository.UserExists"

	sql, args, err := r.qb.
		Select("1").
		Prefix("SELECT EXISTS (").
		From(TableUsers).
		Where(sq.Eq{"oauth_id": oauthID}).
		Suffix(")").
		ToSql()
	if err != nil {
		return false, psql.ErrCreateQuery(op, err)
	}

	var exists bool
	err = r.client.QueryRow(ctx, sql, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, psql.ErrDoQuery(op, err)
	}

	return exists, nil
}

func (r *userRepository) GetUserByOauthID(ctx context.Context, oauthID string) (entity.User, error) {
	logrus.WithField("oauth_id", oauthID).Trace("getting user")
	const op string = "userRepository.GetUser"

	sql, args, err := r.qb.
		Select(
			"id",
			"first_name",
			"last_name",
			"email",
		).
		From(TableUsers).
		Where(sq.Eq{"oauth_id": oauthID}).
		ToSql()
	if err != nil {
		return entity.User{}, psql.ErrCreateQuery(op, err)
	}

	var user entity.User
	err = r.client.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, nil
		}
		return entity.User{}, psql.ErrDoQuery(op, err)
	}

	return user, nil
}

func (r *userRepository) GetUserByUUID(ctx context.Context, userID string) (entity.User, error) {
	logrus.WithField("user_id", userID).Trace("getting user")
	const op string = "userRepository.GetUser"

	sql, args, err := r.qb.
		Select(
			"id",
			"first_name",
			"last_name",
			"email",
		).
		From(TableUsers).
		Where(sq.Eq{"id": userID}).
		ToSql()
	if err != nil {
		return entity.User{}, psql.ErrCreateQuery(op, err)
	}

	var user entity.User
	err = r.client.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, nil
		}
		return entity.User{}, psql.ErrDoQuery(op, err)
	}

	return user, nil

}
