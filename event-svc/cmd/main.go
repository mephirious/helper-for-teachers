package main

import (
	"database/sql"
	"log"
	"net"

	"event-svc/internal/adapters/inbound/grpc"
	"event-svc/internal/adapters/outbound/repository/postgres"
	"event-svc/internal/app/lesson"
	"event-svc/internal/app/schedule"
	"event-svc/internal/app/task"

	"google.golang.org/grpc"
)

func main() {
	// Initialize dependencies
	db := initDB()

	// Repositories
	lessonRepo := postgres.NewLessonRepository(db)
	taskRepo := postgres.NewTaskRepository(db)
	scheduleRepo := postgres.NewScheduleRepository(db)

	// Use cases
	lessonUC := lesson.NewLessonUseCase(lessonRepo)
	taskUC := task.NewTaskUseCase(taskRepo)
	scheduleUC := schedule.NewScheduleUseCase(scheduleRepo)

	// gRPC Server
	grpcServer := grpc.NewServer()
	eventService := grpc.NewEventServiceServer(lessonUC, taskUC, scheduleUC)
	eventsv1.RegisterEventServiceServer(grpcServer, eventService)

	// Start server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Server started on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initDB() *sql.DB {
	// Initialize database connection
	// return sql.Open("postgres", "your-connection-string")
	return nil
}
