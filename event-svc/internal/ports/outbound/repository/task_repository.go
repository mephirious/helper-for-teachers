package repository

import (
	"context"
	"event-svc/internal/domain/model"
	"time"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	GetByID(ctx context.Context, id string) (*model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id string) error
	ListByFilter(ctx context.Context, filter TaskFilter) ([]*model.Task, error)
}

type TaskFilter struct {
	GroupID  *string
	CourseID *string
	Type     *model.TaskType
	Status   *model.TaskStatus
	DateFrom *time.Time
	DateTo   *time.Time
}
