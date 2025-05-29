package repository

import (
	"context"
	"event-svc/internal/domain/model"
	"time"
)

// LessonRepository defines the interface for lesson persistence operations
type LessonRepository interface {
	// Create persists a new lesson and returns its ID
	Create(ctx context.Context, lesson *model.Lesson) (string, error)

	// GetByID retrieves a lesson by its ID
	GetByID(ctx context.Context, id string) (*model.Lesson, error)

	// Update modifies an existing lesson
	Update(ctx context.Context, lesson *model.Lesson) error

	// Delete removes a lesson by its ID
	Delete(ctx context.Context, id string) error

	// GetAll retrieves all lessons
	GetAll(ctx context.Context) ([]*model.Lesson, error)
}

type LessonFilter struct {
	GroupID  *string
	CourseID *string
	DateFrom *time.Time
	DateTo   *time.Time
	Status   *model.LessonStatus
	IsOnline *bool
}
