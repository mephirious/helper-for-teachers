package usecase

import (
	"context"
	"fmt"
	"time"

	redis "github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/redis/cache"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type taskUseCase struct {
	taskRepo  repository.TaskRepository
	taskCache redis.TaskCache
}

func NewTaskUseCase(repo repository.TaskRepository, cache redis.TaskCache) TaskUseCase {
	return &taskUseCase{
		taskRepo:  repo,
		taskCache: cache,
	}
}

func (uc *taskUseCase) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	task.ID = primitive.NewObjectID()
	task.CreatedAt = time.Now()

	if err := uc.taskRepo.CreateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	_ = uc.taskCache.Set(ctx, *task)
	return task, nil
}

func (uc *taskUseCase) GetTaskByID(ctx context.Context, id primitive.ObjectID) (*domain.Task, error) {
	task, err := uc.taskCache.Get(ctx, id.Hex())
	if err == nil {
		return &task, nil
	}

	taskPtr, err := uc.taskRepo.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if taskPtr == nil {
		return nil, fmt.Errorf("task not found")
	}

	_ = uc.taskCache.Set(ctx, *taskPtr)
	return taskPtr, nil
}

func (uc *taskUseCase) GetTasksByExamID(ctx context.Context, examID primitive.ObjectID) ([]domain.Task, error) {
	return uc.taskRepo.GetTasksByExamID(ctx, examID)
}

func (uc *taskUseCase) GetAllTasks(ctx context.Context) ([]domain.Task, error) {
	tasks, err := uc.taskCache.GetAll(ctx)
	if err == nil && len(tasks) > 0 {
		return tasks, nil
	}

	tasks, err = uc.taskRepo.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}

	_ = uc.taskCache.SetMany(ctx, tasks)
	return tasks, nil
}

func (uc *taskUseCase) DeleteTask(ctx context.Context, id primitive.ObjectID) error {
	if err := uc.taskRepo.DeleteTask(ctx, id); err != nil {
		return err
	}
	_ = uc.taskCache.Delete(ctx, id.Hex())
	return nil
}

func (uc *taskUseCase) UpdateTask(ctx context.Context, task *domain.Task) error {
	if err := uc.taskRepo.UpdateTask(ctx, task); err != nil {
		return err
	}
	_ = uc.taskCache.Set(ctx, *task)
	return nil
}
