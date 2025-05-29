package schedule

import (
	"context"
	"event-svc/internal/domain/model"
	"event-svc/internal/ports/outbound/repository"
	"fmt"
)

type ScheduleUseCase struct {
	repo repository.ScheduleRepository
}

func NewScheduleUseCase(repo repository.ScheduleRepository) *ScheduleUseCase {
	return &ScheduleUseCase{repo: repo}
}

// Lesson Schedule CRUD Operations

func (uc *ScheduleUseCase) CreateLessonSchedule(ctx context.Context, schedule *model.LessonSchedule) (*model.LessonSchedule, error) {
	if err := schedule.Validate(); err != nil {
		return nil, fmt.Errorf("lesson schedule validation error: %w", err)
	}

	id, err := uc.repo.CreateLessonSchedule(ctx, schedule)
	if err != nil {
		return nil, fmt.Errorf("failed to create lesson schedule: %w", err)
	}

	createdSchedule, err := uc.repo.GetLessonSchedule(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get created lesson schedule: %w", err)
	}

	return createdSchedule, nil
}

func (uc *ScheduleUseCase) GetLessonSchedule(ctx context.Context, id string) (*model.LessonSchedule, error) {
	if id == "" {
		return nil, model.ErrInvalidID
	}

	schedule, err := uc.repo.GetLessonSchedule(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get lesson schedule: %w", err)
	}

	return schedule, nil
}

func (uc *ScheduleUseCase) UpdateLessonSchedule(ctx context.Context, schedule *model.LessonSchedule) (*model.LessonSchedule, error) {
	if err := schedule.Validate(); err != nil {
		return nil, fmt.Errorf("lesson schedule validation error: %w", err)
	}

	if err := uc.repo.UpdateLessonSchedule(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to update lesson schedule: %w", err)
	}

	updatedSchedule, err := uc.repo.GetLessonSchedule(ctx, schedule.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated lesson schedule: %w", err)
	}

	return updatedSchedule, nil
}

func (uc *ScheduleUseCase) DeleteLessonSchedule(ctx context.Context, id string) error {
	if id == "" {
		return model.ErrInvalidID
	}

	return uc.repo.DeleteLessonSchedule(ctx, id)
}

func (uc *ScheduleUseCase) ListLessonSchedules(ctx context.Context) ([]*model.LessonSchedule, error) {
	schedules, err := uc.repo.ListLessonSchedules(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list lesson schedules: %w", err)
	}

	return schedules, nil
}

// Task Schedule CRUD Operations

func (uc *ScheduleUseCase) CreateTaskSchedule(ctx context.Context, schedule *model.TaskSchedule) (*model.TaskSchedule, error) {
	if err := schedule.Validate(); err != nil {
		return nil, fmt.Errorf("task schedule validation error: %w", err)
	}

	id, err := uc.repo.CreateTaskSchedule(ctx, schedule)
	if err != nil {
		return nil, fmt.Errorf("failed to create task schedule: %w", err)
	}

	createdSchedule, err := uc.repo.GetTaskSchedule(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get created task schedule: %w", err)
	}

	return createdSchedule, nil
}

func (uc *ScheduleUseCase) GetTaskSchedule(ctx context.Context, id string) (*model.TaskSchedule, error) {
	if id == "" {
		return nil, model.ErrInvalidID
	}

	schedule, err := uc.repo.GetTaskSchedule(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task schedule: %w", err)
	}

	return schedule, nil
}

func (uc *ScheduleUseCase) UpdateTaskSchedule(ctx context.Context, schedule *model.TaskSchedule) (*model.TaskSchedule, error) {
	if err := schedule.Validate(); err != nil {
		return nil, fmt.Errorf("task schedule validation error: %w", err)
	}

	if err := uc.repo.UpdateTaskSchedule(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to update task schedule: %w", err)
	}

	updatedSchedule, err := uc.repo.GetTaskSchedule(ctx, schedule.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated task schedule: %w", err)
	}

	return updatedSchedule, nil
}

func (uc *ScheduleUseCase) DeleteTaskSchedule(ctx context.Context, id string) error {
	if id == "" {
		return model.ErrInvalidID
	}

	return uc.repo.DeleteTaskSchedule(ctx, id)
}

func (uc *ScheduleUseCase) ListTaskSchedules(ctx context.Context) ([]*model.TaskSchedule, error) {
	schedules, err := uc.repo.ListTaskSchedules(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list task schedules: %w", err)
	}

	return schedules, nil
}

// Combined Operations

func (uc *ScheduleUseCase) GetSchedulesForGroup(ctx context.Context, groupID string) (*model.GroupSchedules, error) {
	if groupID == "" {
		return nil, model.ErrInvalidID
	}

	lessonSchedules, err := uc.repo.ListLessonSchedulesByGroup(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lesson schedules for group: %w", err)
	}

	taskSchedules, err := uc.repo.ListTaskSchedulesByGroup(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task schedules for group: %w", err)
	}

	return &model.GroupSchedules{
		LessonSchedules: lessonSchedules,
		TaskSchedules:   taskSchedules,
	}, nil
}
