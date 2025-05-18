package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	"event-svc/internal/adapters/inbound/grpc"
	"event-svc/internal/adapters/outbound/repository/postgres"
	"event-svc/internal/app/lesson"
	"event-svc/internal/app/schedule"
	"event-svc/internal/app/task"

	"google.golang.org/grpc"
)

func main() {
	// Initialize database connection
	db, err := initDB()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

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

func initDB() (*sql.DB, error) {
	connStr := os.Getenv("DB_CONNECTION_STRING")
	if connStr == "" {
		connStr = "postgres://user:password@localhost:5432/events?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
