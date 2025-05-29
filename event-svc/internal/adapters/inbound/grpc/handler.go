package grpc

import (
	"context"
	"errors"
	"event-svc/internal/app/lesson"
	"event-svc/internal/app/schedule"
	"event-svc/internal/app/task"
	"event-svc/internal/domain/model"
	"event-svc/internal/ports/outbound/repository"

	eventsv1 "github.com/suyundykovv/margulan-protos/gen/go/events/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventServiceHandler struct {
	eventsv1.UnimplementedEventServiceServer
	lessonUC   *lesson.LessonUseCase
	scheduleUC *schedule.ScheduleUseCase
	taskUC     *task.TaskUseCase
}

func NewEventServiceHandler(
	lessonRepo repository.LessonRepository,
	scheduleRepo repository.ScheduleRepository,
	taskRepo repository.TaskRepository,
) *EventServiceHandler {
	return &EventServiceHandler{
		lessonUC:   lesson.NewLessonUseCase(lessonRepo),
		scheduleUC: schedule.NewScheduleUseCase(scheduleRepo),
		taskUC:     task.NewTaskUseCase(taskRepo),
	}
}

// Lesson handlers
func (h *EventServiceHandler) CreateLesson(ctx context.Context, req *eventsv1.CreateLessonRequest) (*eventsv1.Lesson, error) {
	domainLesson := ConvertCreateLessonRequestToDomain(req)

	createdLesson, err := h.lessonUC.CreateLesson(ctx, domainLesson)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create lesson: %v", err)
	}

	return ConvertDomainLessonToProto(createdLesson), nil
}

func (h *EventServiceHandler) GetLesson(ctx context.Context, req *eventsv1.GetLessonRequest) (*eventsv1.Lesson, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "lesson ID is required")
	}

	lesson, err := h.lessonUC.GetLesson(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "lesson not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get lesson: %v", err)
	}

	return ConvertDomainLessonToProto(lesson), nil
}

// Similar handlers for UpdateLesson, DeleteLesson, ListLessons

// Schedule handlers
func (h *EventServiceHandler) CreateLessonSchedule(ctx context.Context, req *eventsv1.CreateLessonScheduleRequest) (*eventsv1.LessonSchedule, error) {
	domainSchedule := ConvertCreateScheduleRequestToDomain(req)

	createdSchedule, err := h.scheduleUC.CreateLessonSchedule(ctx, domainSchedule)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create schedule: %v", err)
	}

	return ConvertDomainScheduleToProto(createdSchedule), nil
}

// Similar handlers for other schedule operations

// Task handlers
func (h *EventServiceHandler) CreateTask(ctx context.Context, req *eventsv1.CreateTaskRequest) (*eventsv1.Task, error) {
	domainTask := ConvertCreateTaskRequestToDomain(req)

	createdTask, err := h.taskUC.CreateTask(ctx, domainTask)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create task: %v", err)
	}

	return ConvertDomainTaskToProto(createdTask), nil
}

// Similar handlers for other task operations
