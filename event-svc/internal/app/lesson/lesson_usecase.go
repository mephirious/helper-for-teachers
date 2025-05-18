package lesson

import (
	"context"
	"errors"
	"event-svc/internal/domain/model"
	"event-svc/internal/ports/outbound/repository"
	"time"
)

type LessonUseCase struct {
	repo repository.LessonRepository
}

func NewLessonUseCase(repo repository.LessonRepository) *LessonUseCase {
	return &LessonUseCase{repo: repo}
}

func (uc *LessonUseCase) CreateLesson(ctx context.Context, lesson *model.Lesson) (*model.Lesson, error) {
	if err := validateLesson(lesson); err != nil {
		return nil, err
	}

	lesson.ID = generateID()
	lesson.CreatedAt = time.Now()
	lesson.UpdatedAt = time.Now()

	if err := uc.repo.Create(ctx, lesson); err != nil {
		return nil, err
	}

	return lesson, nil
}

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
	return nil
}

func generateID() string {
	// Implement your ID generation logic
	return "generated-id"
}
