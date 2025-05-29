package repository

import (
	"context"
	"event-svc/internal/domain/model"
)

type ScheduleRepository interface {
	// Lesson Schedule operations
	CreateLessonSchedule(ctx context.Context, schedule *model.LessonSchedule) (string, error)
	GetLessonSchedule(ctx context.Context, id string) (*model.LessonSchedule, error)
	UpdateLessonSchedule(ctx context.Context, schedule *model.LessonSchedule) error
	DeleteLessonSchedule(ctx context.Context, id string) error
	ListLessonSchedules(ctx context.Context) ([]*model.LessonSchedule, error)
	ListLessonSchedulesByGroup(ctx context.Context, groupID string) ([]*model.LessonSchedule, error)

	// Task Schedule operations
	CreateTaskSchedule(ctx context.Context, schedule *model.TaskSchedule) (string, error)
	GetTaskSchedule(ctx context.Context, id string) (*model.TaskSchedule, error)
	UpdateTaskSchedule(ctx context.Context, schedule *model.TaskSchedule) error
	DeleteTaskSchedule(ctx context.Context, id string) error
	ListTaskSchedules(ctx context.Context) ([]*model.TaskSchedule, error)
	ListTaskSchedulesByGroup(ctx context.Context, groupID string) ([]*model.TaskSchedule, error)
}
