package task

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/yourproject/event-svc/internal/domain/model"
	"github.com/yourproject/event-svc/internal/ports/outbound/repository"
)

type TaskUseCase struct {
	repo           repository.TaskRepository
	eventPublisher ports.EventPublisher
	notifier       ports.Notifier
}

func NewTaskUseCase(
	repo repository.TaskRepository,
	eventPublisher ports.EventPublisher,
	notifier ports.Notifier,
) *TaskUseCase {
	return &TaskUseCase{
		repo:           repo,
		eventPublisher: eventPublisher,
		notifier:       notifier,
	}
}

func (uc *TaskUseCase) CreateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	if err := validateTask(task); err != nil {
		return nil, err
	}

	task.ID = generateID()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	if err := uc.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	// Publish event
	if err := uc.eventPublisher.PublishTaskCreated(ctx, task); err != nil {
		// Log error but don't fail the operation
		log.Printf("failed to publish task created event: %v", err)
	}

	// Send notifications
	if err := uc.notifier.NotifyTaskAssigned(ctx, task); err != nil {
		log.Printf("failed to send task notifications: %v", err)
	}

	return task, nil
}

func (uc *TaskUseCase) GetTask(ctx context.Context, id string) (*model.Task, error) {
	if id == "" {
		return nil, errors.New("task ID is required")
	}
	return uc.repo.GetByID(ctx, id)
}

func (uc *TaskUseCase) UpdateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	if task.ID == "" {
		return nil, errors.New("task ID is required")
	}

	existing, err := uc.repo.GetByID(ctx, task.ID)
	if err != nil {
		return nil, err
	}

	// Preserve immutable fields
	task.CreatedAt = existing.CreatedAt
	task.UpdatedAt = time.Now()

	if err := validateTask(task); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	// Publish update event if status changed
	if existing.Status != task.Status {
		if err := uc.eventPublisher.PublishTaskStatusChanged(ctx, existing, task); err != nil {
			log.Printf("failed to publish task status changed event: %v", err)
		}
	}

	return task, nil
}

func (uc *TaskUseCase) DeleteTask(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("task ID is required")
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *TaskUseCase) ListTasks(ctx context.Context, filter repository.TaskFilter) ([]*model.Task, error) {
	return uc.repo.ListByFilter(ctx, filter)
}

func (uc *TaskUseCase) GradeTask(ctx context.Context, taskID string, score int32) (*model.Task, error) {
	if taskID == "" {
		return nil, errors.New("task ID is required")
	}

	task, err := uc.repo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if task.MaxScore == nil || *task.MaxScore == 0 {
		return nil, errors.New("task is not gradable")
	}

	if score < 0 || score > *task.MaxScore {
		return nil, errors.New("invalid score")
	}

	task.Status = model.TaskGraded
	task.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	// Publish grading event
	if err := uc.eventPublisher.PublishTaskGraded(ctx, task); err != nil {
		log.Printf("failed to publish task graded event: %v", err)
	}

	// Send grade notification
	if err := uc.notifier.NotifyTaskGraded(ctx, task); err != nil {
		log.Printf("failed to send grade notification: %v", err)
	}

	return task, nil
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
	if task.Type < model.TaskExam || task.Type > model.TaskProject {
		return errors.New("invalid task type")
	}
	return nil
}

func generateID() string {
	// Implement your ID generation logic
	return "generated-id"
}
