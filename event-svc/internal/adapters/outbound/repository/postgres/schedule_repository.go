package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"event-svc/internal/domain/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ScheduleRepository struct {
	db *sql.DB
}

func NewScheduleRepository(db *sql.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

var (
	ErrScheduleNotFound = errors.New("schedule not found")
)

// Lesson Schedule operations

func (r *ScheduleRepository) CreateLessonSchedule(ctx context.Context, schedule *model.LessonSchedule) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	schedule.ID = uuid.New().String()
	schedule.CreatedAt = time.Now()
	schedule.UpdatedAt = time.Now()

	query := `
		INSERT INTO lesson_schedules (
			id, group_id, title, valid_from, valid_to, is_active, 
			course_id, lesson_ids, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err = tx.ExecContext(ctx, query,
		schedule.ID,
		schedule.GroupID,
		schedule.Title,
		schedule.ValidFrom,
		schedule.ValidTo,
		schedule.IsActive,
		schedule.CourseID,
		pq.Array(schedule.LessonIDs),
		schedule.CreatedAt,
		schedule.UpdatedAt,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create lesson schedule: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return schedule.ID, nil
}

func (r *ScheduleRepository) GetLessonSchedule(ctx context.Context, id string) (*model.LessonSchedule, error) {
	query := `
		SELECT * FROM lesson_schedules WHERE id = $1
	`

	var schedule model.LessonSchedule
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&schedule.ID,
		&schedule.GroupID,
		&schedule.Title,
		&schedule.ValidFrom,
		&schedule.ValidTo,
		&schedule.IsActive,
		&schedule.CourseID,
		pq.Array(&schedule.LessonIDs),
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrScheduleNotFound
		}
		return nil, fmt.Errorf("failed to get lesson schedule: %w", err)
	}

	return &schedule, nil
}

func (r *ScheduleRepository) UpdateLessonSchedule(ctx context.Context, schedule *model.LessonSchedule) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	schedule.UpdatedAt = time.Now()

	query := `
		UPDATE lesson_schedules SET
			title = $1,
			valid_to = $2,
			is_active = $3,
			lesson_ids = $4,
			updated_at = $5
		WHERE id = $6
	`

	result, err := tx.ExecContext(ctx, query,
		schedule.Title,
		schedule.ValidTo,
		schedule.IsActive,
		pq.Array(schedule.LessonIDs),
		schedule.UpdatedAt,
		schedule.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update lesson schedule: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrScheduleNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *ScheduleRepository) DeleteLessonSchedule(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `DELETE FROM lesson_schedules WHERE id = $1`

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete lesson schedule: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrScheduleNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *ScheduleRepository) ListLessonSchedules(ctx context.Context) ([]*model.LessonSchedule, error) {
	query := `
		SELECT * FROM lesson_schedules
		ORDER BY valid_from ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list lesson schedules: %w", err)
	}
	defer rows.Close()

	var schedules []*model.LessonSchedule
	for rows.Next() {
		var schedule model.LessonSchedule
		err := rows.Scan(
			&schedule.ID,
			&schedule.GroupID,
			&schedule.Title,
			&schedule.ValidFrom,
			&schedule.ValidTo,
			&schedule.IsActive,
			&schedule.CourseID,
			pq.Array(&schedule.LessonIDs),
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lesson schedule: %w", err)
		}
		schedules = append(schedules, &schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning rows: %w", err)
	}

	return schedules, nil
}

func (r *ScheduleRepository) ListLessonSchedulesByGroup(ctx context.Context, groupID string) ([]*model.LessonSchedule, error) {
	query := `
		SELECT * FROM lesson_schedules
		WHERE group_id = $1
		ORDER BY valid_from ASC
	`

	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list lesson schedules by group: %w", err)
	}
	defer rows.Close()

	var schedules []*model.LessonSchedule
	for rows.Next() {
		var schedule model.LessonSchedule
		err := rows.Scan(
			&schedule.ID,
			&schedule.GroupID,
			&schedule.Title,
			&schedule.ValidFrom,
			&schedule.ValidTo,
			&schedule.IsActive,
			&schedule.CourseID,
			pq.Array(&schedule.LessonIDs),
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lesson schedule: %w", err)
		}
		schedules = append(schedules, &schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning rows: %w", err)
	}

	return schedules, nil
}

// Task Schedule operations

func (r *ScheduleRepository) CreateTaskSchedule(ctx context.Context, schedule *model.TaskSchedule) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	schedule.ID = uuid.New().String()
	schedule.CreatedAt = time.Now()
	schedule.UpdatedAt = time.Now()

	query := `
		INSERT INTO task_schedules (
			id, group_id, title, valid_from, valid_to, is_active, 
			course_id, task_ids, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err = tx.ExecContext(ctx, query,
		schedule.ID,
		schedule.GroupID,
		schedule.Title,
		schedule.ValidFrom,
		schedule.ValidTo,
		schedule.IsActive,
		schedule.CourseID,
		pq.Array(schedule.TaskIDs),
		schedule.CreatedAt,
		schedule.UpdatedAt,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create task schedule: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return schedule.ID, nil
}

func (r *ScheduleRepository) GetTaskSchedule(ctx context.Context, id string) (*model.TaskSchedule, error) {
	query := `
		SELECT * FROM task_schedules WHERE id = $1
	`

	var schedule model.TaskSchedule
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&schedule.ID,
		&schedule.GroupID,
		&schedule.Title,
		&schedule.ValidFrom,
		&schedule.ValidTo,
		&schedule.IsActive,
		&schedule.CourseID,
		pq.Array(&schedule.TaskIDs),
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrScheduleNotFound
		}
		return nil, fmt.Errorf("failed to get task schedule: %w", err)
	}

	return &schedule, nil
}

func (r *ScheduleRepository) UpdateTaskSchedule(ctx context.Context, schedule *model.TaskSchedule) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	schedule.UpdatedAt = time.Now()

	query := `
		UPDATE task_schedules SET
			title = $1,
			valid_to = $2,
			is_active = $3,
			task_ids = $4,
			updated_at = $5
		WHERE id = $6
	`

	result, err := tx.ExecContext(ctx, query,
		schedule.Title,
		schedule.ValidTo,
		schedule.IsActive,
		pq.Array(schedule.TaskIDs),
		schedule.UpdatedAt,
		schedule.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task schedule: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrScheduleNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *ScheduleRepository) DeleteTaskSchedule(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `DELETE FROM task_schedules WHERE id = $1`

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task schedule: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrScheduleNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *ScheduleRepository) ListTaskSchedules(ctx context.Context) ([]*model.TaskSchedule, error) {
	query := `
		SELECT * FROM task_schedules
		ORDER BY valid_from ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list task schedules: %w", err)
	}
	defer rows.Close()

	var schedules []*model.TaskSchedule
	for rows.Next() {
		var schedule model.TaskSchedule
		err := rows.Scan(
			&schedule.ID,
			&schedule.GroupID,
			&schedule.Title,
			&schedule.ValidFrom,
			&schedule.ValidTo,
			&schedule.IsActive,
			&schedule.CourseID,
			pq.Array(&schedule.TaskIDs),
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task schedule: %w", err)
		}
		schedules = append(schedules, &schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning rows: %w", err)
	}

	return schedules, nil
}

func (r *ScheduleRepository) ListTaskSchedulesByGroup(ctx context.Context, groupID string) ([]*model.TaskSchedule, error) {
	query := `
		SELECT * FROM task_schedules
		WHERE group_id = $1
		ORDER BY valid_from ASC
	`

	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list task schedules by group: %w", err)
	}
	defer rows.Close()

	var schedules []*model.TaskSchedule
	for rows.Next() {
		var schedule model.TaskSchedule
		err := rows.Scan(
			&schedule.ID,
			&schedule.GroupID,
			&schedule.Title,
			&schedule.ValidFrom,
			&schedule.ValidTo,
			&schedule.IsActive,
			&schedule.CourseID,
			pq.Array(&schedule.TaskIDs),
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task schedule: %w", err)
		}
		schedules = append(schedules, &schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning rows: %w", err)
	}

	return schedules, nil
}
