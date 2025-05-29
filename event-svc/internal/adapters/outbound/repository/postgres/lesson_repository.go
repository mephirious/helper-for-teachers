package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"event-svc/internal/domain/model"

	"github.com/google/uuid"
)

type LessonRepository struct {
	db *sql.DB
}

func NewLessonRepository(db *sql.DB) *LessonRepository {
	return &LessonRepository{db: db}
}

func (r *LessonRepository) Create(ctx context.Context, lesson *model.Lesson) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	lesson.ID = uuid.New().String()
	lesson.CreatedAt = time.Now()
	lesson.UpdatedAt = time.Now()

	query := `INSERT INTO lessons (id, title, start_time, end_time, group_id, course_id, 
		status, meeting_url, classroom, is_online, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err = tx.ExecContext(ctx, query,
		lesson.ID,
		lesson.Title,
		lesson.StartTime,
		lesson.EndTime,
		lesson.GroupID,
		lesson.CourseID,
		lesson.Status,
		lesson.MeetingURL,
		lesson.Classroom,
		lesson.IsOnline,
		lesson.CreatedAt,
		lesson.UpdatedAt,
	)

	if err != nil {
		return "", fmt.Errorf("failed to create lesson: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return lesson.ID, nil
}

func (r *LessonRepository) GetByID(ctx context.Context, id string) (*model.Lesson, error) {
	query := `
        SELECT 
            id, 
            title, 
            start_time, 
            end_time, 
            group_id, 
            course_id, 
            status, 
            meeting_url, 
            classroom, 
            is_online, 
            created_at, 
            updated_at
        FROM lessons 
        WHERE id = $1
    `

	row := r.db.QueryRowContext(ctx, query, id)
	var lesson model.Lesson
	err := row.Scan(
		&lesson.ID,
		&lesson.Title,
		&lesson.StartTime,
		&lesson.EndTime,
		&lesson.GroupID,
		&lesson.CourseID,
		&lesson.Status,
		&lesson.MeetingURL,
		&lesson.Classroom,
		&lesson.IsOnline,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get lesson: %w", err)
	}

	return &lesson, nil
}

func (r *LessonRepository) GetAll(ctx context.Context) ([]*model.Lesson, error) {
	baseQuery := `
        SELECT 
            id, 
            title, 
            start_time, 
            end_time, 
            group_id, 
            course_id, 
            status, 
            meeting_url, 
            classroom, 
            is_online, 
            created_at, 
            updated_at
        FROM lessons
    `

	var conditions []string
	var args []interface{}

	finalQuery := baseQuery
	if len(conditions) > 0 {
		finalQuery += " WHERE " + strings.Join(conditions, " AND ")
	}
	finalQuery += " ORDER BY start_time ASC"

	rows, err := r.db.QueryContext(ctx, finalQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get lessons: %w", err)
	}
	defer rows.Close()

	var lessons []*model.Lesson
	for rows.Next() {
		var lesson model.Lesson
		err := rows.Scan(
			&lesson.ID,
			&lesson.Title,
			&lesson.StartTime,
			&lesson.EndTime,
			&lesson.GroupID,
			&lesson.CourseID,
			&lesson.Status,
			&lesson.MeetingURL,
			&lesson.Classroom,
			&lesson.IsOnline,
			&lesson.CreatedAt,
			&lesson.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lesson: %w", err)
		}
		lessons = append(lessons, &lesson)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning lessons: %w", err)
	}

	return lessons, nil
}

func (r *LessonRepository) Update(ctx context.Context, lesson *model.Lesson) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	lesson.UpdatedAt = time.Now()

	query := `
		UPDATE lessons SET
			title = $1,
			start_time = $2,
			end_time = $3,
			status = $4,
			meeting_url = $5,
			classroom = $6,
			is_online = $7,
			updated_at = $8
		WHERE id = $9
	`

	result, err := tx.ExecContext(ctx, query,
		lesson.Title,
		lesson.StartTime,
		lesson.EndTime,
		lesson.Status,
		lesson.MeetingURL,
		lesson.Classroom,
		lesson.IsOnline,
		lesson.UpdatedAt,
		lesson.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update lesson: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return model.ErrNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *LessonRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `DELETE FROM lessons WHERE id = $1`

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete lesson: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return model.ErrNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
