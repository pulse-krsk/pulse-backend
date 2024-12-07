package postgresql

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	psql "github.com/kurochkinivan/pulskrsk/pkg/postgresql"
	"github.com/sirupsen/logrus"
)

type usersFavoriteTypesRepository struct {
	client psql.PosgreSQLClient
	qb     sq.StatementBuilderType
}

func NewUsersFavoriteTypesRepository(client *pgxpool.Pool) *usersFavoriteTypesRepository {
	return &usersFavoriteTypesRepository{
		client: client,
		qb:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *usersFavoriteTypesRepository) AddUserFavouriteTypes(ctx context.Context, userID string, typeIDs []string) error {
	logrus.WithField("user_id", userID).Debug("adding favorite event types")
	const op string = "usersFavoriteTypesRepository.AddUserFavouriteTypes"

	if len(typeIDs) == 0 {
		return nil
	}

	query := r.qb.Insert(TableUsersFavoriteTypes).
		Columns("user_id", "type_id")
	for _, typeID := range typeIDs {
		query = query.Values(userID, typeID)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return psql.ErrCreateQuery(op, err)
	}
	commTag, err := r.client.Exec(ctx, sql, args...)
	if err != nil {
		return psql.ErrExec(op, err)
	}

	if commTag.RowsAffected() == 0 {
		return psql.NoRowsAffected
	}

	return nil
}

func (r *usersFavoriteTypesRepository) DeleteAllUserFavouriteTypes(ctx context.Context, userID string) error {
	logrus.WithFields(logrus.Fields{"user_id": userID}).Debug("deleting favorite event types")
	const op string = "usersFavoriteTypesRepository.DeleteAllUserFavouriteTypes"

	sql, args, err := r.qb.
		Delete(TableUsersFavoriteTypes).
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return psql.ErrCreateQuery(op, err)
	}

	_, err = r.client.Exec(ctx, sql, args...)
	if err != nil {
		return psql.ErrExec(op, err)
	}

	return nil
}

func (r *usersFavoriteTypesRepository) GetAllUserFavouriteTypesID(ctx context.Context, userID string) ([]string, error) {
	logrus.WithFields(logrus.Fields{"user_id": userID}).Debug("getting favorite event types")
	const op string = "usersFavoriteTypesRepository.GetAllUserFavouriteTypesID"

	sql, args, err := r.qb.
		Select("type_id").
		From(TableUsersFavoriteTypes).
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, psql.ErrCreateQuery(op, err)
	}

	rows, err := r.client.Query(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, psql.ErrDoQuery(op, err)
	}

	var typeIDs []string
	for rows.Next() {
		var typeID string
		err = rows.Scan(&typeID)
		if err != nil {
			return nil, psql.ErrScan(op, err)
		}
		typeIDs = append(typeIDs, typeID)
	}

	return typeIDs, nil
}
