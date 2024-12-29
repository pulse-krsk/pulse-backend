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

type EventTypeRepository struct {
	client psql.PosgreSQLClient
	qb     sq.StatementBuilderType
}

func NewEventTypeRepository(client *pgxpool.Pool) *EventTypeRepository {
	return &EventTypeRepository{
		client: client,
		qb:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *EventTypeRepository) GetEventTypeByType(ctx context.Context, eventType string) (entity.EventType, error) {
	logrus.WithField("event_type", eventType).Debug("getting event by type")
	const op string = "EventTypeRepository.GetEventByType"

	query, args, err := r.qb.
		Select(
			"id",
			"type",
		).
		From(TableEventTypes).
		Where(sq.Eq{"type": eventType}).
		ToSql()
	if err != nil {
		return entity.EventType{}, psql.ErrCreateQuery(op, err)
	}

	var newEventType entity.EventType
	err = r.client.QueryRow(ctx, query, args...).Scan(&newEventType.ID, &newEventType.Type)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.EventType{}, nil
		}
		return entity.EventType{}, psql.ErrDoQuery(op, err)
	}

	return newEventType, nil
}

func (r *EventTypeRepository) GetEventTypeByID(ctx context.Context, typeID string) (entity.EventType, error) {
	logrus.WithField("event_type_id", typeID).Debug("getting event by id")
	const op string = "EventTypeRepository.GetEventTypeByID"

	sql, args, err := r.qb.
		Select(
			"id",
			"type",
		).
		From(TableEventTypes).
		Where(sq.Eq{"id": typeID}).
		ToSql()
	if err != nil {
		return entity.EventType{}, psql.ErrCreateQuery(op, err)
	}

	var eventType entity.EventType
	err = r.client.QueryRow(ctx, sql, args...).Scan(&eventType.ID, &eventType.Type)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.EventType{}, nil
		}
		return entity.EventType{}, psql.ErrDoQuery(op, err)
	}

	return eventType, nil
}
