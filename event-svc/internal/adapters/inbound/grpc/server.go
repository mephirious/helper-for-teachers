package grpc

import (
	"event-svc/internal/app/lesson"
	"event-svc/internal/app/schedule"
	"event-svc/internal/app/task"
	"fmt"
	"log"
	"net"

	eventsv1 "github.com/suyundykovv/margulan-protos/gen/go/events/v1"

	"google.golang.org/grpc"
)

// Server wraps the gRPC server and its dependencies
type Server struct {
	grpcServer   *grpc.Server
	lessonUC     *lesson.LessonUseCase
	taskUC       *task.TaskUseCase
	scheduleUC   *schedule.ScheduleUseCase
	eventService *EventServiceServer
}

// NewServer creates a new gRPC server instance with all dependencies
func NewServer(
	lessonUC *lesson.LessonUseCase,
	taskUC *task.TaskUseCase,
	scheduleUC *schedule.ScheduleUseCase,
) *Server {
	// Create gRPC server with options (add interceptors if needed)
	grpcServer := grpc.NewServer()

	// Initialize the event service handler
	eventService := NewEventServiceServer(lessonUC, taskUC, scheduleUC)

	return &Server{
		grpcServer:   grpcServer,
		lessonUC:     lessonUC,
		taskUC:       taskUC,
		scheduleUC:   scheduleUC,
		eventService: eventService,
	}
}

// Start runs the gRPC server on the specified address
func (s *Server) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Register the event service
	eventsv1.RegisterEventServiceServer(s.grpcServer, s.eventService)

	log.Printf("Starting gRPC server on %s", addr)
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

// GracefulStop shuts down the gRPC server gracefully
func (s *Server) GracefulStop() {
	log.Println("Shutting down gRPC server gracefully...")
	s.grpcServer.GracefulStop()
}

// EventServiceServer implements the gRPC service
type EventServiceServer struct {
	eventsv1.UnimplementedEventServiceServer
	lessonUC   *lesson.LessonUseCase
	taskUC     *task.TaskUseCase
	scheduleUC *schedule.ScheduleUseCase
}

// NewEventServiceServer creates a new EventServiceServer instance
func NewEventServiceServer(
	lessonUC *lesson.LessonUseCase,
	taskUC *task.TaskUseCase,
	scheduleUC *schedule.ScheduleUseCase,
) *EventServiceServer {
	return &EventServiceServer{
		lessonUC:   lessonUC,
		taskUC:     taskUC,
		scheduleUC: scheduleUC,
	}
}
