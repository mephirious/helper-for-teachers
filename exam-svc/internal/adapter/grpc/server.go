package service

import (
	"fmt"
	"net"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/config"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/grpc/handler"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/usecase"
	pb "github.com/mephirious/helper-for-teachers/services/exam-svc/proto"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	Cfg      config.GRPCServer
	server   *grpc.Server
	addr     string
	listener net.Listener
}

func NewGRPCServer(cfg config.Config, taskUC usecase.TaskUseCase, questionUC usecase.QuestionUseCase, examUC usecase.ExamUseCase) (*GRPCServer, error) {
	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Server.GRPCServer.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	examHandler := handler.NewExamHandler(taskUC, questionUC, examUC)

	pb.RegisterExamServiceServer(s, examHandler)
	reflection.Register(s)

	return &GRPCServer{
		Cfg:      cfg.Server.GRPCServer,
		server:   s,
		addr:     addr,
		listener: lis,
	}, nil
}

func (s *GRPCServer) Run() error {
	fmt.Printf("gRPC server running on %s\n", s.addr)
	return s.server.Serve(s.listener)
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
	fmt.Println("gRPC server stopped gracefully")
}
