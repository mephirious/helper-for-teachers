package service

import (
	"context"
	"errors"
	"time"

	"event-svc/internal/domain/model"
)

// LessonService defines the interface for lesson operations
type LessonService interface {
	CreateLesson(ctx context.Context, lesson *model.Lesson) (*model.Lesson, error)
	GetLesson(ctx context.Context, id string) (*model.Lesson, error)
	UpdateLesson(ctx context.Context, lesson *model.Lesson) (*model.Lesson, error)
	DeleteLesson(ctx context.Context, id string) error
	ListLessons(ctx context.Context, filter LessonFilter) ([]*model.Lesson, error)
}

// LessonFilter defines filtering options for lessons
type LessonFilter struct {
	GroupID  *string
	CourseID *string
	DateFrom *time.Time
	DateTo   *time.Time
	Status   *model.LessonStatus
	IsOnline *bool
}

// lessonServiceImpl implements LessonService
type lessonServiceImpl struct {
	repo model.LessonRepository
}

// NewLessonService creates a new lesson service
func NewLessonService(repo model.LessonRepository) LessonService {
	return &lessonServiceImpl{repo: repo}
}

func (s *lessonServiceImpl) CreateLesson(ctx context.Context, lesson *model.Lesson) (*model.Lesson, error) {
	// Validate lesson
	if err := validateLesson(lesson); err != nil {
		return nil, err
	}

	// Set default values
	lesson.ID = generateID()
	lesson.CreatedAt = time.Now()
	lesson.UpdatedAt = time.Now()

	// Persist to repository
	if err := s.repo.Create(ctx, lesson); err != nil {
		return nil, err
	}

	return lesson, nil
}

func (s *lessonServiceImpl) GetLesson(ctx context.Context, id string) (*model.Lesson, error) {
	if id == "" {
		return nil, errors.New("lesson ID is required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *lessonServiceImpl) UpdateLesson(ctx context.Context, lesson *model.Lesson) (*model.Lesson, error) {
	if lesson.ID == "" {
		return nil, errors.New("lesson ID is required")
	}

	// Validate lesson
	if err := validateLesson(lesson); err != nil {
		return nil, err
	}

	// Set update timestamp
	lesson.UpdatedAt = time.Now()

	// Update in repository
	if err := s.repo.Update(ctx, lesson); err != nil {
		return nil, err
	}

	return lesson, nil
}

func (s *lessonServiceImpl) DeleteLesson(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("lesson ID is required")
	}
	return s.repo.Delete(ctx, id)
}

func (s *lessonServiceImpl) ListLessons(ctx context.Context, filter LessonFilter) ([]*model.Lesson, error) {
	return s.repo.ListByFilter(ctx, convertToRepoFilter(filter))
}

// Helper functions
func validateLesson(lesson *model.Lesson) error {
	if lesson.Title == "" {
		return errors.New("title is required")
	}
	if lesson.StartTime.IsZero() {
		return errors.New("start time is required")
	}
	if lesson.EndTime.IsZero() {
		return errors.New("end time is required")
	}
	if lesson.StartTime.After(lesson.EndTime) {
		return errors.New("start time must be before end time")
	}
	if lesson.GroupID == "" {
		return errors.New("group ID is required")
	}
	if lesson.CourseID == "" {
		return errors.New("course ID is required")
	}
	return nil
}

func generateID() string {
	// Implement your actual ID generation logic
	return "generated-id"
}

func convertToRepoFilter(filter LessonFilter) model.LessonFilter {
	return model.LessonFilter{
		GroupID:  filter.GroupID,
		CourseID: filter.CourseID,
		DateFrom: filter.DateFrom,
		DateTo:   filter.DateTo,
		Status:   filter.Status,
		IsOnline: filter.IsOnline,
	}
}
