package main

import (
	"event-svc/internal/adapters/inbound/grpc"
	postgres "event-svc/internal/adapters/outbound/repository/postgres"
	"event-svc/internal/app/lesson"
	"event-svc/internal/app/schedule"
	"event-svc/internal/app/task"
	"event-svc/pkg/config"
	"log"
)

func main() {
	// Initialize database connection
	db, _ := config.InitDB()
	defer db.Close()

	// Initialize repositories
	lessonRepo := postgres.NewLessonRepository(db)
	taskRepo := postgres.NewTaskRepository(db)
	scheduleRepo := postgres.NewScheduleRepository(db)

	// Initialize use cases
	lessonUC := lesson.NewLessonUseCase(lessonRepo)
	taskUC := task.NewTaskUseCase(taskRepo)
	scheduleUC := schedule.NewScheduleUseCase(scheduleRepo)

	// Create and start gRPC server
	grpcServer := grpc.NewServer(lessonUC, taskUC, scheduleUC)

	if err := grpcServer.Start(":50051"); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
