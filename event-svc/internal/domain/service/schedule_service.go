package service

import (
	"context"
	"errors"
	"time"

	"event-svc/internal/domain/model"
)

// ScheduleService defines the interface for schedule operations
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

// LessonScheduleFilter defines filtering options for lesson schedules
type LessonScheduleFilter struct {
	GroupID  *string
	CourseID *string
	IsActive *bool
	ActiveAt *time.Time
}

// TaskScheduleFilter defines filtering options for task schedules
type TaskScheduleFilter struct {
	GroupID  *string
	CourseID *string
	IsActive *bool
	ActiveAt *time.Time
}

// scheduleServiceImpl implements ScheduleService
type scheduleServiceImpl struct {
	repo model.ScheduleRepository
}

// NewScheduleService creates a new schedule service
func NewScheduleService(repo model.ScheduleRepository) ScheduleService {
	return &scheduleServiceImpl{repo: repo}
}

// Lesson Schedule methods
func (s *scheduleServiceImpl) CreateLessonSchedule(ctx context.Context, schedule *model.LessonSchedule) (*model.LessonSchedule, error) {
	if err := validateLessonSchedule(schedule); err != nil {
		return nil, err
	}

	schedule.ID = generateID()
	schedule.CreatedAt = time.Now()
	schedule.UpdatedAt = time.Now()

	if err := s.repo.CreateLessonSchedule(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *scheduleServiceImpl) GetLessonSchedule(ctx context.Context, id string) (*model.LessonSchedule, error) {
	if id == "" {
		return nil, errors.New("schedule ID is required")
	}
	return s.repo.GetLessonSchedule(ctx, id)
}

func (s *scheduleServiceImpl) UpdateLessonSchedule(ctx context.Context, schedule *model.LessonSchedule) (*model.LessonSchedule, error) {
	if schedule.ID == "" {
		return nil, errors.New("schedule ID is required")
	}

	if err := validateLessonSchedule(schedule); err != nil {
		return nil, err
	}

	schedule.UpdatedAt = time.Now()

	if err := s.repo.UpdateLessonSchedule(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *scheduleServiceImpl) DeleteLessonSchedule(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("schedule ID is required")
	}
	return s.repo.DeleteLessonSchedule(ctx, id)
}

func (s *scheduleServiceImpl) ListLessonSchedules(ctx context.Context, filter LessonScheduleFilter) ([]*model.LessonSchedule, error) {
	return s.repo.ListLessonSchedules(ctx, model.LessonScheduleFilter{
		GroupID:  filter.GroupID,
		CourseID: filter.CourseID,
		IsActive: filter.IsActive,
		ActiveAt: filter.ActiveAt,
	})
}

// Task Schedule methods
func (s *scheduleServiceImpl) CreateTaskSchedule(ctx context.Context, schedule *model.TaskSchedule) (*model.TaskSchedule, error) {
	if err := validateTaskSchedule(schedule); err != nil {
		return nil, err
	}

	schedule.ID = generateID()
	schedule.CreatedAt = time.Now()
	schedule.UpdatedAt = time.Now()

	if err := s.repo.CreateTaskSchedule(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *scheduleServiceImpl) GetTaskSchedule(ctx context.Context, id string) (*model.TaskSchedule, error) {
	if id == "" {
		return nil, errors.New("schedule ID is required")
	}
	return s.repo.GetTaskSchedule(ctx, id)
}

func (s *scheduleServiceImpl) UpdateTaskSchedule(ctx context.Context, schedule *model.TaskSchedule) (*model.TaskSchedule, error) {
	if schedule.ID == "" {
		return nil, errors.New("schedule ID is required")
	}

	if err := validateTaskSchedule(schedule); err != nil {
		return nil, err
	}

	schedule.UpdatedAt = time.Now()

	if err := s.repo.UpdateTaskSchedule(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *scheduleServiceImpl) DeleteTaskSchedule(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("schedule ID is required")
	}
	return s.repo.DeleteTaskSchedule(ctx, id)
}

func (s *scheduleServiceImpl) ListTaskSchedules(ctx context.Context, filter TaskScheduleFilter) ([]*model.TaskSchedule, error) {
	return s.repo.ListTaskSchedules(ctx, model.TaskScheduleFilter{
		GroupID:  filter.GroupID,
		CourseID: filter.CourseID,
		IsActive: filter.IsActive,
		ActiveAt: filter.ActiveAt,
	})
}

// Validation helpers
func validateLessonSchedule(schedule *model.LessonSchedule) error {
	if schedule.Title == "" {
		return errors.New("title is required")
	}
	if schedule.GroupID == "" {
		return errors.New("group ID is required")
	}
	if schedule.CourseID == "" {
		return errors.New("course ID is required")
	}
	if schedule.ValidFrom.IsZero() {
		return errors.New("valid from date is required")
	}
	if schedule.ValidTo.IsZero() {
		return errors.New("valid to date is required")
	}
	if schedule.ValidFrom.After(schedule.ValidTo) {
		return errors.New("valid from must be before valid to")
	}
	return nil
}

func validateTaskSchedule(schedule *model.TaskSchedule) error {
	if schedule.Title == "" {
		return errors.New("title is required")
	}
	if schedule.GroupID == "" {
		return errors.New("group ID is required")
	}
	if schedule.CourseID == "" {
		return errors.New("course ID is required")
	}
	if schedule.ValidFrom.IsZero() {
		return errors.New("valid from date is required")
	}
	if schedule.ValidTo.IsZero() {
		return errors.New("valid to date is required")
	}
	if schedule.ValidFrom.After(schedule.ValidTo) {
		return errors.New("valid from must be before valid to")
	}
	return nil
}
