package grpc

import (
	"event-svc/internal/ports"

	eventsv1 "github.com/suyundykovv/margulan-protos/gen/go/events/v1"

	"google.golang.org/grpc"
)

func NewServer(
	lessonService ports.LessonService,
	taskService ports.TaskService,
	scheduleService ports.ScheduleService,
) *grpc.Server {
	srv := grpc.NewServer()
	handler := NewEventServiceServer(lessonService, taskService, scheduleService)
	eventsv1.RegisterEventServiceServer(srv, handler)
	return srv
}
