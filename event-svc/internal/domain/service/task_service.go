package service

import (
	"context"
	"time"

	"event-svc/internal/domain/model"
)

// TaskService defines the interface for task operations
type TaskService interface {
	CreateTask(ctx context.Context, task *model.Task) (*model.Task, error)
	GetTask(ctx context.Context, id string) (*model.Task, error)
	UpdateTask(ctx context.Context, task *model.Task) (*model.Task, error)
	DeleteTask(ctx context.Context, id string) error
	ListTasks(ctx context.Context, filter TaskFilter) ([]*model.Task, error)
}

// TaskFilter defines filtering options for tasks
type TaskFilter struct {
	GroupID  *string
	CourseID *string
	Type     *model.TaskType
	Status   *model.TaskStatus
	DateFrom *time.Time
	DateTo   *time.Time
}
