package service

import (
	"context"
	"time"

	"event-svc/internal/domain/model"
)

type ScheduleService interface {
	CreateLessonSchedule(ctx context.Context, schedule *model.LessonSchedule) (*model.LessonSchedule, error)
	GetLessonSchedule(ctx context.Context, id string) (*model.LessonSchedule, error)
	UpdateLessonSchedule(ctx context.Context, schedule *model.LessonSchedule) (*model.LessonSchedule, error)
	DeleteLessonSchedule(ctx context.Context, id string) error
	ListLessonSchedules(ctx context.Context, filter LessonScheduleFilter) ([]*model.LessonSchedule, error)

	CreateTaskSchedule(ctx context.Context, schedule *model.TaskSchedule) (*model.TaskSchedule, error)
	GetTaskSchedule(ctx context.Context, id string) (*model.TaskSchedule, error)
	UpdateTaskSchedule(ctx context.Context, schedule *model.TaskSchedule) (*model.TaskSchedule, error)
	DeleteTaskSchedule(ctx context.Context, id string) error
	ListTaskSchedules(ctx context.Context, filter TaskScheduleFilter) ([]*model.TaskSchedule, error)
}

type LessonScheduleFilter struct {
	GroupID  *string
	CourseID *string
	IsActive *bool
	ActiveAt *time.Time
}

type TaskScheduleFilter struct {
	GroupID  *string
	CourseID *string
	IsActive *bool
	ActiveAt *time.Time
}
