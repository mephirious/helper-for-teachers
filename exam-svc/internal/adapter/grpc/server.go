package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/grpc/handler"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/usecase"
	pb "github.com/mephirious/helper-for-teachers/services/exam-svc/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	logger     *logrus.Logger
}

func NewServer(taskUseCase usecase.TaskUseCase, questionUseCase usecase.QuestionUseCase, examUseCase usecase.ExamUseCase) *Server {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	examHandler := handler.NewExamHandler(taskUseCase, questionUseCase, examUseCase)
	grpcServer := grpc.NewServer()
	pb.RegisterExamServiceServer(grpcServer, examHandler)

	return &Server{
		grpcServer: grpcServer,
		logger:     logger,
	}
}

func (s *Server) Start(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		s.logger.Errorf("failed to listen on port %s: %v", port, err)
		return err
	}

	s.logger.Infof("Starting gRPC server on port %s", port)
	go func() {
		if err := s.grpcServer.Serve(listener); err != nil {
			s.logger.Errorf("failed to serve gRPC server: %v", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.logger.Warn("Shutdown timeout, forcing stop")
		s.grpcServer.Stop()
		return ctx.Err()
	case <-stopped:
		s.logger.Info("gRPC server stopped gracefully")
		return nil
	}
}

func (s *Server) Run(port string) error {
	if err := s.Start(port); err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	s.logger.Info("Received shutdown signal")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.Shutdown(ctx)
}
