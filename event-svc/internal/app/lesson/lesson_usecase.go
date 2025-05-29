package lesson

import (
	"context"
	"event-svc/internal/domain/model"
	"event-svc/internal/ports/outbound/repository"
	"fmt"
)

type LessonUseCase struct {
	repo repository.LessonRepository
}

func NewLessonUseCase(repo repository.LessonRepository) *LessonUseCase {
	return &LessonUseCase{repo: repo}
}

func (uc *LessonUseCase) CreateLesson(ctx context.Context, lesson *model.Lesson) (*model.Lesson, error) {
	if err := lesson.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	id, err := uc.repo.Create(ctx, lesson)
	if err != nil {
		return nil, fmt.Errorf("failed to create lesson: %w", err)
	}

	createdLesson, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get created lesson: %w", err)
	}

	return createdLesson, nil
}

func (uc *LessonUseCase) GetLesson(ctx context.Context, id string) (*model.Lesson, error) {
	if id == "" {
		return nil, model.ErrInvalidID
	}

	lesson, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get lesson: %w", err)
	}

	return lesson, nil
}

func (uc *LessonUseCase) UpdateLesson(ctx context.Context, lesson *model.Lesson) (*model.Lesson, error) {
	if err := lesson.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if err := uc.repo.Update(ctx, lesson); err != nil {
		return nil, fmt.Errorf("failed to update lesson: %w", err)
	}

	updatedLesson, err := uc.repo.GetByID(ctx, lesson.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated lesson: %w", err)
	}

	return updatedLesson, nil
}

func (uc *LessonUseCase) DeleteLesson(ctx context.Context, id string) error {
	if id == "" {
		return model.ErrInvalidID
	}

	return uc.repo.Delete(ctx, id)
}

func (uc *LessonUseCase) ListLessons(ctx context.Context) ([]*model.Lesson, error) {
	lessons, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list lessons: %w", err)
	}

	return lessons, nil
}
