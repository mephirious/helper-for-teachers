package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/core/domain"
)

type PostgresInstructorsRepository struct {
	db *sql.DB
}

func NewInstructorsRepository(db *sql.DB) *PostgresInstructorsRepository {
	return &PostgresInstructorsRepository{db: db}
}

func (r *PostgresInstructorsRepository) Create(ctx context.Context, instr *domain.CourseInstructor) error {
	const op = "PostgresInstructorsRepository.Create"
	return runInTx(ctx, r.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO course_instructors (
				id, course_id, user_id, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5)`,
			instr.ID,
			instr.CourseID,
			instr.UserID,
			instr.CreatedAt,
			instr.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
}

func (r *PostgresInstructorsRepository) Delete(ctx context.Context, courseID, userID uuid.UUID) error {
	const op = "PostgresInstructorsRepository.Delete"
	return runInTx(ctx, r.db, func(tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx,
			`DELETE FROM course_instructors WHERE course_id = $1 AND user_id = $2`,
			courseID,
			userID,
		)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if affected == 0 {
			fmt.Println("Here!", courseID, userID)
			return domain.ErrCourseInstructorNotFound
		}
		return nil
	})
}

func (r *PostgresInstructorsRepository) ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*domain.CourseInstructor, error) {
	const op = "PostgresInstructorsRepository.ListByCourse"
	var instructors []*domain.CourseInstructor

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, course_id, user_id, created_at, updated_at
		FROM course_instructors WHERE course_id = $1`,
		courseID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var i domain.CourseInstructor
		if err := rows.Scan(&i.ID, &i.CourseID, &i.UserID, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		instructors = append(instructors, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return instructors, nil
}
