package postgresql

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kurochkinivan/pulskrsk/internal/entity"
	psql "github.com/kurochkinivan/pulskrsk/pkg/postgresql"
)

type roleRepository struct {
	client psql.PosgreSQLClient
	qb     sq.StatementBuilderType
}

func NewRoleRepository(client *pgxpool.Pool) *roleRepository {
	return &roleRepository{
		client: client,
		qb:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *roleRepository) CreateRole(ctx context.Context, role entity.Role) error {
	const op string = "roleRepository.CreateRole"

	sql, args, err := r.qb.
		Insert(TableRole).
		Columns(
			"user_id",
			"role",
		).
		Values(
			role.UserID,
			role.RoleName,
		).
		ToSql()
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

func (r *roleRepository) GetRolesByUserID(ctx context.Context, userID string) ([]entity.Role, error) {
	const op string = "roleRepository.GetRolesByUserID"
	sql, args, err := r.qb.
		Select(
			"id",
			"role",
			"user_id",
		).
		From(TableRole).
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, psql.ErrCreateQuery(op, err)
	}

	var roles []entity.Role
	rows, err := r.client.Query(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, psql.ErrDoQuery(op, err)
	}

	for rows.Next() {
		var role entity.Role
		err = rows.Scan(&role.ID, &role.RoleName, &role.UserID)
		if err != nil {
			return nil, psql.ErrDoQuery(op, err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}
