package service

import (
	"context"
	"time"

	"event-svc/internal/domain/model"
)

type LessonService interface {
	CreateLesson(ctx context.Context, lesson *model.Lesson) (*model.Lesson, error)
	GetLesson(ctx context.Context, id string) (*model.Lesson, error)
	UpdateLesson(ctx context.Context, lesson *model.Lesson) (*model.Lesson, error)
	DeleteLesson(ctx context.Context, id string) error
	ListLessons(ctx context.Context, filter LessonFilter) ([]*model.Lesson, error)
}

type LessonFilter struct {
	GroupID  *string
	CourseID *string
	DateFrom *time.Time
	DateTo   *time.Time
	Status   *model.LessonStatus
	IsOnline *bool
}
