package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/core/domain"
)

type PostgresCourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *PostgresCourseRepository {
	return &PostgresCourseRepository{db: db}
}

func (r *PostgresCourseRepository) Create(ctx context.Context, course *domain.Course) error {
	const op = "PostgresCourseRepository.Create"
	return runInTx(ctx, r.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO courses (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)`,
			course.ID, course.Name, course.CreatedAt, course.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
}

func (r *PostgresCourseRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Course, error) {
	const op = "PostgresCourseRepository.GetByID"
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, created_at, updated_at FROM courses WHERE id = $1`, id,
	)

	var course domain.Course
	if err := row.Scan(&course.ID, &course.Name, &course.CreatedAt, &course.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCourseNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &course, nil
}

func (r *PostgresCourseRepository) Update(ctx context.Context, id uuid.UUID, updateFn func(*domain.Course) (bool, error)) error {
	const op = "PostgresCourseRepository.Update"
	return runInTx(ctx, r.db, func(tx *sql.Tx) error {
		var course domain.Course
		err := tx.QueryRowContext(ctx,
			`SELECT id, name, created_at, updated_at FROM courses WHERE id = $1 FOR UPDATE`, id,
		).Scan(&course.ID, &course.Name, &course.CreatedAt, &course.UpdatedAt)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return domain.ErrCourseNotFound
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		updated, err := updateFn(&course)
		if err != nil {
			return err
		}
		if !updated {
			return domain.ErrNotUpdated
		}

		course.UpdatedAt = time.Now().UTC()
		_, err = tx.ExecContext(ctx,
			`UPDATE courses SET name = $1, updated_at = $2 WHERE id = $3`,
			course.Name, course.UpdatedAt, course.ID,
		)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
}

func (r *PostgresCourseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "PostgresCourseRepository.Delete"
	return runInTx(ctx, r.db, func(tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `DELETE FROM courses WHERE id = $1`, id)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		n, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if n == 0 {
			return domain.ErrCourseNotFound
		}
		return nil
	})
}

func (r *PostgresCourseRepository) List(ctx context.Context) ([]*domain.Course, error) {
	const op = "PostgresCourseRepository.List"
	var courses []*domain.Course

	err := runInTx(ctx, r.db, func(tx *sql.Tx) error {
		rows, err := tx.QueryContext(ctx, `SELECT id, name, created_at, updated_at FROM courses`)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		defer rows.Close()

		for rows.Next() {
			var course domain.Course
			if err := rows.Scan(&course.ID, &course.Name, &course.CreatedAt, &course.UpdatedAt); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			courses = append(courses, &course)
		}
		return rows.Err()
	})

	return courses, err
}
