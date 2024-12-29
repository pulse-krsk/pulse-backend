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

	refreshToken, _ := uuid.NewUUID()
	sql, args, err := r.qb.
		Insert(TableRefreshSessions).
		Columns(
			"refresh_token",
			"user_id",
			"issued_at",
			"expiration",
		).
		Values(
			refreshToken.String(),
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

func (r *refreshSessionsRepository) DeleteRefreshSessionByToken(ctx context.Context, refreshToken string) error {
	logrus.WithField("refresh_token", refreshToken).Trace("deleting refresh session by refresh token")
	const op string = "refreshSessionsRepository.DeleteRefreshSessionByToken"

	sql, args, err := r.qb.
		Delete(TableRefreshSessions).
		Where(sq.Eq{"refresh_token": refreshToken}).
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
