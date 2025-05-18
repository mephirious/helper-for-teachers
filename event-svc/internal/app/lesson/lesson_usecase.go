package schedule

import (
	"context"
	"errors"
	"time"

	"event-svc/internal/domain/model"
	"event-svc/internal/ports/outbound/repository"
)

type ScheduleUseCase struct {
	lessonRepo     repository.LessonRepository
	taskRepo       repository.TaskRepository
	scheduleRepo   repository.ScheduleRepository
	eventPublisher ports.EventPublisher
}

func NewScheduleUseCase(
	lessonRepo repository.LessonRepository,
	taskRepo repository.TaskRepository,
	scheduleRepo repository.ScheduleRepository,
	eventPublisher ports.EventPublisher,
) *ScheduleUseCase {
	return &ScheduleUseCase{
		lessonRepo:     lessonRepo,
		taskRepo:       taskRepo,
		scheduleRepo:   scheduleRepo,
		eventPublisher: eventPublisher,
	}
}

// Lesson Schedule methods
func (uc *ScheduleUseCase) CreateLessonSchedule(ctx context.Context, schedule *model.LessonSchedule) (*model.LessonSchedule, error) {
	if err := validateLessonSchedule(schedule); err != nil {
		return nil, err
	}

	schedule.ID = generateID()
	schedule.CreatedAt = time.Now()
	schedule.UpdatedAt = time.Now()

	if err := uc.scheduleRepo.CreateLessonSchedule(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (uc *ScheduleUseCase) AddLessonToSchedule(ctx context.Context, scheduleID string, lessonID string) error {
	if scheduleID == "" || lessonID == "" {
		return errors.New("schedule ID and lesson ID are required")
	}

	// Verify lesson exists
	if _, err := uc.lessonRepo.GetByID(ctx, lessonID); err != nil {
		return err
	}

	schedule, err := uc.scheduleRepo.GetLessonSchedule(ctx, scheduleID)
	if err != nil {
		return err
	}

	// Check if lesson already exists in schedule
	for _, id := range schedule.LessonIDs {
		if id == lessonID {
			return nil // already exists
		}
	}

	schedule.LessonIDs = append(schedule.LessonIDs, lessonID)
	schedule.UpdatedAt = time.Now()

	return uc.scheduleRepo.UpdateLessonSchedule(ctx, schedule)
}

func (uc *ScheduleUseCase) RemoveLessonFromSchedule(ctx context.Context, scheduleID string, lessonID string) error {
	if scheduleID == "" || lessonID == "" {
		return errors.New("schedule ID and lesson ID are required")
	}

	schedule, err := uc.scheduleRepo.GetLessonSchedule(ctx, scheduleID)
	if err != nil {
		return err
	}

	var newLessonIDs []string
	for _, id := range schedule.LessonIDs {
		if id != lessonID {
			newLessonIDs = append(newLessonIDs, id)
		}
	}

	if len(newLessonIDs) == len(schedule.LessonIDs) {
		return nil // lesson not in schedule
	}

	schedule.LessonIDs = newLessonIDs
	schedule.UpdatedAt = time.Now()

	return uc.scheduleRepo.UpdateLessonSchedule(ctx, schedule)
}

func (uc *ScheduleUseCase) GetLessonSchedule(ctx context.Context, id string) (*model.LessonSchedule, error) {
	if id == "" {
		return nil, errors.New("schedule ID is required")
	}
	return uc.scheduleRepo.GetLessonSchedule(ctx, id)
}

func (uc *ScheduleUseCase) ListLessonSchedules(ctx context.Context, filter repository.LessonScheduleFilter) ([]*model.LessonSchedule, error) {
	return uc.scheduleRepo.ListLessonSchedules(ctx, filter)
}

// Task Schedule methods
func (uc *ScheduleUseCase) CreateTaskSchedule(ctx context.Context, schedule *model.TaskSchedule) (*model.TaskSchedule, error) {
	if err := validateTaskSchedule(schedule); err != nil {
		return nil, err
	}

	schedule.ID = generateID()
	schedule.CreatedAt = time.Now()
	schedule.UpdatedAt = time.Now()

	if err := uc.scheduleRepo.CreateTaskSchedule(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (uc *ScheduleUseCase) AddTaskToSchedule(ctx context.Context, scheduleID string, taskID string) error {
	if scheduleID == "" || taskID == "" {
		return errors.New("schedule ID and task ID are required")
	}

	// Verify task exists
	if _, err := uc.taskRepo.GetByID(ctx, taskID); err != nil {
		return err
	}

	schedule, err := uc.scheduleRepo.GetTaskSchedule(ctx, scheduleID)
	if err != nil {
		return err
	}

	// Check if task already exists in schedule
	for _, id := range schedule.TaskIDs {
		if id == taskID {
			return nil // already exists
		}
	}

	schedule.TaskIDs = append(schedule.TaskIDs, taskID)
	schedule.UpdatedAt = time.Now()

	return uc.scheduleRepo.UpdateTaskSchedule(ctx, schedule)
}

func (uc *ScheduleUseCase) RemoveTaskFromSchedule(ctx context.Context, scheduleID string, taskID string) error {
	if scheduleID == "" || taskID == "" {
		return errors.New("schedule ID and task ID are required")
	}

	schedule, err := uc.scheduleRepo.GetTaskSchedule(ctx, scheduleID)
	if err != nil {
		return err
	}

	var newTaskIDs []string
	for _, id := range schedule.TaskIDs {
		if id != taskID {
			newTaskIDs = append(newTaskIDs, id)
		}
	}

	if len(newTaskIDs) == len(schedule.TaskIDs) {
		return nil // task not in schedule
	}

	schedule.TaskIDs = newTaskIDs
	schedule.UpdatedAt = time.Now()

	return uc.scheduleRepo.UpdateTaskSchedule(ctx, schedule)
}

func (uc *ScheduleUseCase) GetTaskSchedule(ctx context.Context, id string) (*model.TaskSchedule, error) {
	if id == "" {
		return nil, errors.New("schedule ID is required")
	}
	return uc.scheduleRepo.GetTaskSchedule(ctx, id)
}

func (uc *ScheduleUseCase) ListTaskSchedules(ctx context.Context, filter repository.TaskScheduleFilter) ([]*model.TaskSchedule, error) {
	return uc.scheduleRepo.ListTaskSchedules(ctx, filter)
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

func generateID() string {
	// Implement your ID generation logic
	return "generated-id"
}
