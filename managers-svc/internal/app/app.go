package app

import (
	"database/sql"
	"log"
	"log/slog"
	"time"

	"github.com/mephirious/helper-for-teachers/managers-svc/config"
	smtp "github.com/mephirious/helper-for-teachers/managers-svc/internal/SMTP"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/cache"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/core/service"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/grpc"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/nats"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/repository"

	_ "github.com/lib/pq"
)

type APIServer struct {
	cfg    *config.Config
	logger *slog.Logger
}

func NewAPIServer(cfg *config.Config, logger *slog.Logger) *APIServer {
	return &APIServer{
		cfg:    cfg,
		logger: logger,
	}
}
func (s *APIServer) Run() error {
	connStr := s.cfg.Database.MakeConnectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("could not open database: %v", err)
	}
	courseRepo := repository.NewCourseRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	instructorsRepo := repository.NewInstructorsRepository(db)
	membersRepo := repository.NewMembersRepository(db)
	cache := cache.NewCacheRepository(membersRepo)

	publisher, err := nats.NewEventPublisher(s.cfg.NATS.URL, s.logger)
	if err != nil {
		log.Fatalf("Failed to create NATS publisher: %v", err)
	}
	defer publisher.Close()

	SMTP := smtp.NewClient(s.cfg.SMTP.Host, int(s.cfg.SMTP.Port), s.cfg.SMTP.Email, s.cfg.SMTP.Password)
	err = SMTP.Send([]string{"2004gusak@gmail.com"}, "Service started", "The service is started, this are good news!")
	if err != nil {
		s.logger.Error("SMTP ERROR", "error", err.Error())
	}
	service := service.NewManagersService(time.Now, s.logger, courseRepo, groupRepo, cache, instructorsRepo, publisher, SMTP)

	return grpc.StartGRPCServer(s.cfg.Server.Port, service, s.logger)
}
