package postgresql

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kurochkinivan/pulskrsk/internal/entity"
	psql "github.com/kurochkinivan/pulskrsk/pkg/postgresql"
	"github.com/sirupsen/logrus"
)

type refreshSessionsRepository struct {
	client psql.PosgreSQLClient
	qb     sq.StatementBuilderType
}

func NewRefreshSessionsRepository(client *pgxpool.Pool) *refreshSessionsRepository {
	return &refreshSessionsRepository{
		client: client,
		qb:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// TODO: think about no rows affected error

func (r *refreshSessionsRepository) CreateRefreshSession(ctx context.Context, refreshSession entity.RefreshSession) (string, error) {
	logrus.WithField("user_id", refreshSession.UserID).Trace("creating refresh session")
	const op string = "refreshSessionsRepository.CreateRefreshSession"

	sql, args, err := r.qb.
		Insert(TableRefreshSessions).
		Columns(
			"user_id",
			"issued_at",
			"expiration",
		).Values(
		refreshSession.UserID,
		refreshSession.IssuedAt,
		refreshSession.Expiration,
	).Suffix("RETURNING refresh_token").
		ToSql()
	if err != nil {
		return "", psql.ErrCreateQuery(op, err)
	}

	var refresh string
	err = r.client.QueryRow(ctx, sql, args...).Scan(&refresh)
	if err != nil {
		return "", psql.ErrDoQuery(op, err)
	}

	return refresh, nil
}

func (r *refreshSessionsRepository) GetRefreshSession(ctx context.Context, refreshToken string) (entity.RefreshSession, error) {
	logrus.WithField("refresh_token", refreshToken).Trace("getting refresh session")
	const op string = "refreshSessionsRepository.GetRefreshSession"

	sql, args, err := r.qb.
		Select(
			"user_id",
			"issued_at",
			"expiration",
		).
		From(TableRefreshSessions).
		Where(sq.Eq{"refresh_token": refreshToken}).
		ToSql()
	if err != nil {
		return entity.RefreshSession{}, psql.ErrCreateQuery(op, err)
	}

	var refreshSession entity.RefreshSession
	err = r.client.QueryRow(ctx, sql, args...).Scan(&refreshSession.UserID, &refreshSession.IssuedAt, &refreshSession.Expiration)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.RefreshSession{}, nil
		}
		return entity.RefreshSession{}, psql.ErrDoQuery(op, err)
	}

	return refreshSession, nil
}

func (r *refreshSessionsRepository) DeleteRefreshSessionsByUserID(ctx context.Context, userID string) error {
	return nil
}
func (r *refreshSessionsRepository) DeleteRefreshSessionByToken(ctx context.Context, refreshToken string) error {
	return nil
}