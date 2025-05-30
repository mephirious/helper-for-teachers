package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/core/domain"
)

type PostgresGroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *PostgresGroupRepository {
	return &PostgresGroupRepository{db: db}
}

func (r *PostgresGroupRepository) Create(ctx context.Context, group *domain.Group) error {
	const op = "PostgresGroupRepository.Create"
	return runInTx(ctx, r.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO groups (id, course_id, name, expire_at, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, group.ID, group.CourseID, group.Name, group.ExpireAt, group.CreatedAt, group.UpdatedAt)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
}

func (r *PostgresGroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
	const op = "PostgresGroupRepository.GetByID"
	var g domain.Group
	err := r.db.QueryRowContext(ctx, `
		SELECT id, course_id, name, expire_at, created_at, updated_at
		FROM groups
		WHERE id = $1
	`, id).Scan(&g.ID, &g.CourseID, &g.Name, &g.ExpireAt, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrGroupNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &g, nil
}

func (r *PostgresGroupRepository) Update(ctx context.Context, id uuid.UUID, updateFn func(*domain.Group) (bool, error)) error {
	const op = "PostgresGroupRepository.Update"
	return runInTx(ctx, r.db, func(tx *sql.Tx) error {
		var group domain.Group
		err := tx.QueryRowContext(ctx,
			`SELECT id, course_id, name, expire_at, created_at, updated_at
			 FROM groups
			 WHERE id = $1 FOR UPDATE`, id,
		).Scan(
			&group.ID,
			&group.CourseID,
			&group.Name,
			&group.ExpireAt,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return domain.ErrGroupNotFound
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		updated, err := updateFn(&group)
		if err != nil {
			return err
		}
		if !updated {
			return domain.ErrNotUpdated
		}

		_, err = tx.ExecContext(ctx,
			`UPDATE groups SET name = $1, expire_at = $2, updated_at = $3 WHERE id = $4`,
			group.Name, group.ExpireAt, group.UpdatedAt, group.ID,
		)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
}

func (r *PostgresGroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "PostgresGroupRepository.Delete"
	return runInTx(ctx, r.db, func(tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `DELETE FROM groups WHERE id = $1`, id)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return domain.ErrGroupNotFound
		}
		return nil
	})
}

func (r *PostgresGroupRepository) ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*domain.Group, error) {
	const op = "PostgresGroupRepository.ListByCourse"
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, course_id, name, expire_at, created_at, updated_at
		FROM groups
		WHERE course_id = $1
	`, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var groups []*domain.Group
	for rows.Next() {
		var g domain.Group
		err := rows.Scan(&g.ID, &g.CourseID, &g.Name, &g.ExpireAt, &g.CreatedAt, &g.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		groups = append(groups, &g)
	}
	return groups, nil
}
