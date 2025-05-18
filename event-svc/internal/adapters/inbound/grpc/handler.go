package grpc

import (
	"context"

	"event-svc/internal/ports"

	eventsv1 "github.com/suyundykovv/margulan-protos/gen/go/events/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventServiceServer struct {
	eventsv1.UnimplementedEventServiceServer
	lessonService   ports.LessonService
	taskService     ports.TaskService
	scheduleService ports.ScheduleService
}

func NewEventServiceServer(
	lessonService ports.LessonService,
	taskService ports.TaskService,
	scheduleService ports.ScheduleService,
) *EventServiceServer {
	return &EventServiceServer{
		lessonService:   lessonService,
		taskService:     taskService,
		scheduleService: scheduleService,
	}
}

func (s *EventServiceServer) CreateLesson(ctx context.Context, req *eventsv1.CreateLessonRequest) (*eventsv1.Lesson, error) {
	lesson, err := convertCreateLessonRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	created, err := s.lessonService.CreateLesson(ctx, lesson)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertLessonToProto(created), nil
}

// Implement all other gRPC methods following the same pattern
