package task

import (
	"context"
	"event-svc/internal/domain/model"
	"event-svc/internal/ports/outbound/repository"
	"fmt"
)

type TaskUseCase struct {
	repo repository.TaskRepository
}

func NewTaskUseCase(repo repository.TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

func (uc *TaskUseCase) CreateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	if err := task.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	id, err := uc.repo.CreateTask(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	createdTask, err := uc.repo.GetTask(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get created task: %w", err)
	}

	return createdTask, nil
}

func (uc *TaskUseCase) GetTask(ctx context.Context, id string) (*model.Task, error) {
	if id == "" {
		return nil, model.ErrInvalidID
	}

	task, err := uc.repo.GetTask(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

func (uc *TaskUseCase) UpdateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	if err := task.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if err := uc.repo.UpdateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	updatedTask, err := uc.repo.GetTask(ctx, task.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated task: %w", err)
	}

	return updatedTask, nil
}

func (uc *TaskUseCase) DeleteTask(ctx context.Context, id string) error {
	if id == "" {
		return model.ErrInvalidID
	}

	return uc.repo.DeleteTask(ctx, id)
}

func (uc *TaskUseCase) ListTasks(ctx context.Context) ([]*model.Task, error) {
	tasks, err := uc.repo.ListTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return tasks, nil
}

func (uc *TaskUseCase) BatchCreateTasks(ctx context.Context, tasks []*model.Task) ([]*model.Task, error) {
	for _, task := range tasks {
		if err := task.Validate(); err != nil {
			return nil, fmt.Errorf("validation error: %w", err)
		}
	}

	ids, err := uc.repo.BatchCreateTasks(ctx, tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to batch create tasks: %w", err)
	}

	var createdTasks []*model.Task
	for _, id := range ids {
		task, err := uc.repo.GetTask(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get created task: %w", err)
		}
		createdTasks = append(createdTasks, task)
	}

	return createdTasks, nil
}
