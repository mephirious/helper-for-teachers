package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type taskUseCase struct {
	taskRepo repository.TaskRepository
}

func NewTaskUseCase(repo repository.TaskRepository) TaskUseCase {
	return &taskUseCase{
		taskRepo: repo,
	}
}

func (uc *taskUseCase) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	task.ID = primitive.NewObjectID()
	task.CreatedAt = time.Now()

	err := uc.taskRepo.CreateTask(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	return task, nil
}

func (uc *taskUseCase) GetTaskByID(ctx context.Context, id primitive.ObjectID) (*domain.Task, error) {
	task, err := uc.taskRepo.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}
	return task, nil
}

func (uc *taskUseCase) GetTasksByExamID(ctx context.Context, examID primitive.ObjectID) ([]domain.Task, error) {
	return uc.taskRepo.GetTasksByExamID(ctx, examID)
}

func (uc *taskUseCase) GetAllTasks(ctx context.Context) ([]domain.Task, error) {
	return uc.taskRepo.GetAllTasks(ctx)
}

func (uc *taskUseCase) DeleteTask(ctx context.Context, id primitive.ObjectID) error {
	return uc.taskRepo.DeleteTask(ctx, id)
}

func (u *taskUseCase) UpdateTask(ctx context.Context, task *domain.Task) error {
	return u.taskRepo.UpdateTask(ctx, task)
}
