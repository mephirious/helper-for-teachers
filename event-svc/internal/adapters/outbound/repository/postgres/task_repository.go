package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"event-svc/internal/domain/model"

	"github.com/google/uuid"
	"github.com/lib/pq" // or whatever driver you're using
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

var (
	ErrTaskNotFound = errors.New("task not found")
)

// CreateTask creates a new task
func (r *TaskRepository) CreateTask(ctx context.Context, task *model.Task) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	task.ID = uuid.New().String()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	query := `
		INSERT INTO tasks (
			id, title, description, due_date, group_id, course_id, type, status,
			external_resource_id, attachments, max_score, lesson_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`

	_, err = tx.ExecContext(ctx, query,
		task.ID,
		task.Title,
		task.Description,
		task.DueDate,
		task.GroupID,
		task.CourseID,
		task.Type,
		task.Status,
		pq.Array(task.Attachments), // assuming you're using PostgreSQL and need array handling
		task.MaxScore,
		task.LessonID,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create task: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return task.ID, nil
}

// GetTask retrieves a task by ID
func (r *TaskRepository) GetTask(ctx context.Context, id string) (*model.Task, error) {
	query := `
		SELECT * FROM tasks WHERE id = $1
	`

	var task model.Task
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.GroupID,
		&task.CourseID,
		&task.Type,
		&task.Status,
		pq.Array(&task.Attachments),
		&task.MaxScore,
		&task.LessonID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

// UpdateTask updates an existing task
func (r *TaskRepository) UpdateTask(ctx context.Context, task *model.Task) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	task.UpdatedAt = time.Now()

	query := `
		UPDATE tasks SET
			title = $1,
			description = $2,
			due_date = $3,
			type = $4,
			status = $5,
			external_resource_id = $6,
			attachments = $7,
			max_score = $8,
			lesson_id = $9,
			updated_at = $10
		WHERE id = $11
	`

	result, err := tx.ExecContext(ctx, query,
		task.Title,
		task.Description,
		task.DueDate,
		task.Type,
		task.Status,
		pq.Array(task.Attachments),
		task.MaxScore,
		task.LessonID,
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ListTasks retrieves all tasks
func (r *TaskRepository) ListTasks(ctx context.Context) ([]*model.Task, error) {
	query := `
		SELECT * FROM tasks
		ORDER BY due_date ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		var task model.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&task.GroupID,
			&task.CourseID,
			&task.Type,
			&task.Status,
			pq.Array(&task.Attachments),
			&task.MaxScore,
			&task.LessonID,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning rows: %w", err)
	}

	return tasks, nil
}

// DeleteTask deletes a task by ID
func (r *TaskRepository) DeleteTask(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

// BatchCreateTasks creates multiple tasks in a single transaction
func (r *TaskRepository) BatchCreateTasks(ctx context.Context, tasks []*model.Task) ([]string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO tasks (
			id, title, description, due_date, group_id, course_id, type, status,
			external_resource_id, attachments, max_score, lesson_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`

	var ids []string
	now := time.Now()

	for _, task := range tasks {
		task.ID = uuid.New().String()
		task.CreatedAt = now
		task.UpdatedAt = now

		_, err := tx.ExecContext(ctx, query,
			task.ID,
			task.Title,
			task.Description,
			task.DueDate,
			task.GroupID,
			task.CourseID,
			task.Type,
			task.Status,
			pq.Array(task.Attachments),
			task.MaxScore,
			task.LessonID,
			task.CreatedAt,
			task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create task %s: %w", task.Title, err)
		}
		ids = append(ids, task.ID)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ids, nil
}
