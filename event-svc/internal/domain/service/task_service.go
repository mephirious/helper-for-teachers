package service

import (
	"context"
	"errors"
	"time"

	"github.com/yourproject/event-svc/internal/domain/model"
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

// taskServiceImpl implements TaskService
type taskServiceImpl struct {
	repo model.TaskRepository
}

// NewTaskService creates a new task service
func NewTaskService(repo model.TaskRepository) TaskService {
	return &taskServiceImpl{repo: repo}
}

func (s *taskServiceImpl) CreateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	if err := validateTask(task); err != nil {
		return nil, err
	}

	task.ID = generateID()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskServiceImpl) GetTask(ctx context.Context, id string) (*model.Task, error) {
	if id == "" {
		return nil, errors.New("task ID is required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *taskServiceImpl) UpdateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	if task.ID == "" {
		return nil, errors.New("task ID is required")
	}

	if err := validateTask(task); err != nil {
		return nil, err
	}

	task.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskServiceImpl) DeleteTask(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("task ID is required")
	}
	return s.repo.Delete(ctx, id)
}

func (s *taskServiceImpl) ListTasks(ctx context.Context, filter TaskFilter) ([]*model.Task, error) {
	return s.repo.ListByFilter(ctx, model.TaskFilter{
		GroupID:  filter.GroupID,
		CourseID: filter.CourseID,
		Type:     filter.Type,
		Status:   filter.Status,
		DateFrom: filter.DateFrom,
		DateTo:   filter.DateTo,
	})
}

func validateTask(task *model.Task) error {
	if task.Title == "" {
		return errors.New("title is required")
	}
	if task.DueDate.IsZero() {
		return errors.New("due date is required")
	}
	if task.GroupID == "" {
		return errors.New("group ID is required")
	}
	if task.CourseID == "" {
		return errors.New("course ID is required")
	}
	return nil
}
