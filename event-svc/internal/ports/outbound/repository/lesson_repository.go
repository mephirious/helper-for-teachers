package repository

import (
	"context"
	"event-svc/internal/domain/model"
	"time"
)

type LessonRepository interface {
	Create(ctx context.Context, lesson *model.Lesson) error
	GetByID(ctx context.Context, id string) (*model.Lesson, error)
	Update(ctx context.Context, lesson *model.Lesson) error
	Delete(ctx context.Context, id string) error
	ListByFilter(ctx context.Context, filter LessonFilter) ([]*model.Lesson, error)
}

type LessonFilter struct {
	GroupID  *string
	CourseID *string
	DateFrom *time.Time
	DateTo   *time.Time
	Status   *model.LessonStatus
	IsOnline *bool
}
