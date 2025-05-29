package repository

import (
	"context"
	"event-svc/internal/domain/model"
	"time"
)

// TaskRepository defines the interface for task persistence operations
type TaskRepository interface {
	// CreateTask persists a new task and returns its ID
	CreateTask(ctx context.Context, task *model.Task) (string, error)

	// GetTask retrieves a task by its ID
	GetTask(ctx context.Context, id string) (*model.Task, error)

	// UpdateTask modifies an existing task
	UpdateTask(ctx context.Context, task *model.Task) error

	// DeleteTask removes a task by its ID
	DeleteTask(ctx context.Context, id string) error

	// ListTasks retrieves all tasks
	ListTasks(ctx context.Context) ([]*model.Task, error)

	// BatchCreateTasks creates multiple tasks in a single transaction
	BatchCreateTasks(ctx context.Context, tasks []*model.Task) ([]string, error)
}
type TaskFilter struct {
	GroupID  *string
	CourseID *string
	Type     *model.TaskType
	Status   *model.TaskStatus
	DateFrom *time.Time
	DateTo   *time.Time
}
