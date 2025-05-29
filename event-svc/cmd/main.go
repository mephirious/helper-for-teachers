package main

import (
	"log"
	"net"
	"os"

	natss "event-scv/internal/adapters/outbound/nats"
	"event-svc/internal/adapters/inbound/grpc"
	"event-svc/internal/adapters/outbound/repository/postgres"
	schedule "event-svc/internal/app/lesson"
	task "event-svc/internal/app/schedule"
	"event-svc/internal/app/task"
	config "event-svc/pkg/config"

	eventsv1 "github.com/suyundykovv/margulan-protos/gen/go/events/v1"
)

func main() {
	cfg := config.Load()

	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	nc, err := natss.Connect(cfg.NATS.URL, 5, 5)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// Initialize repositories
	lessonRepo := postgres.NewLessonRepository(db)
	taskRepo := postgres.NewTaskRepository(db)
	scheduleRepo := postgres.NewScheduleRepository(db)

	// Initialize use cases
	lessonUC := lesson.NewLessonUseCase(lessonRepo)
	taskUC := task.NewTaskUseCase(taskRepo)
	scheduleUC := schedule.NewScheduleUseCase(scheduleRepo)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	eventService := grpc.NewEventServiceServer(lessonUC, taskUC, scheduleUC)
	eventsv1.RegisterEventServiceServer(grpcServer, eventService)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Server started on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
