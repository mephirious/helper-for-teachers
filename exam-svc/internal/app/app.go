package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/config"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/gemini"
	service "github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/grpc"
	memo "github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/in-memory"
	mail "github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/mailjet"
	observability "github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/metrics"
	natsAdapter "github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/nats"
	cache "github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/redis/cache"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/repository"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/usecase"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/pkg/mongo"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/pkg/nats"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/pkg/redis"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/genai"
)

const serviceName = "exam-service"

type App struct {
	grpcServer *service.GRPCServer
	mongo      *mongo.DB
	redis      *redis.Client
	nats       *nats.Client
	gemini     *genai.Client
	tracer     trace.TracerProvider
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	log.Printf("starting %v service", serviceName)

	tracer, err := observability.InitTracer(ctx, "exam-svc", "jaeger:4318")
	if err != nil {
		log.Printf("failed to initialize tracer: %v", err)
		return nil, fmt.Errorf("tracer: %w", err)
	}
	log.Println("initialized OpenTelemetry tracer")

	log.Println("connecting to MongoDB", "uri", cfg.Mongo.URI)
	mongoClient, err := mongo.NewDB(ctx, cfg.Mongo)
	if err != nil {
		log.Printf("failed to initialize MongoDB client: %v", err)
		return nil, fmt.Errorf("mongo: %w", err)
	}
	log.Println("connected to MongoDB")

	log.Println("connecting to Redis", "addr", cfg.Redis.Addr)
	redisClient, err := redis.NewClient(ctx, cfg.Redis)
	if err != nil {
		log.Printf("failed to initialize Redis client: %v", err)
		return nil, fmt.Errorf("redis: %w", err)
	}
	log.Println("connected to Redis")

	log.Println("connecting to NATS", "url", cfg.NATS.URL)
	natsConn, err := nats.NewClient(cfg.NATS.URL)
	if err != nil {
		log.Printf("failed to initialize NATS client: %v", err)
		return nil, fmt.Errorf("nats: %w", err)
	}
	log.Println("connected to NATS")

	log.Println("initializing Gemini client")
	geminiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: cfg.Gemini.APIKey,
	})
	if err != nil {
		log.Printf("failed to initialize Gemini client: %v", err)
		return nil, fmt.Errorf("gemini: %w", err)
	}
	log.Println("initialized Gemini client")

	_, err = observability.InitMetrics(serviceName)
	if err != nil {
		log.Printf("failed to initialize Prometheus metrics: %v", err)
		return nil, fmt.Errorf("metrics: %w", err)
	}
	log.Println("initialized Prometheus metrics")

	mongoDB, err := mongo.NewDB(ctx, cfg.Mongo)
	if err != nil {
		return nil, fmt.Errorf("mongo: %w", err)
	}

	examRepo := repository.NewExamRepository(mongoDB.Connection, mongoClient.Client)
	questionRepo := repository.NewQuestionRepository(mongoDB.Connection, mongoClient.Client)
	taskRepo := repository.NewTaskRepository(mongoDB.Connection, mongoClient.Client)

	questionCache := cache.NewQuestionCache(redisClient)
	questionCache.Init(ctx, questionRepo)
	taskCache := cache.NewTaskCache(redisClient)
	taskCache.Init(ctx, taskRepo)

	cacheManager := memo.NewCacheManager()
	cacheManager.ExamCache.Init(ctx, examRepo)
	cacheManager.QuestionCache.Init(ctx, questionRepo)
	cacheManager.TaskCache.Init(ctx, taskRepo)

	publisher := natsAdapter.NewExamEventProducer(natsConn, "exam.events")
	geminiAdapter, err := gemini.NewClient(geminiClient, cfg.Gemini.ModelName)
	if err != nil {
		log.Printf("failed to initialize gemini adapter: %v", err)
		return nil, fmt.Errorf("grpc server: %w", err)
	}

	client := mailjet.NewMailjetClient(cfg.Mailjet.API, cfg.Mailjet.KEY)
	mailer := mail.NewMailjetClient(client, cfg.Mailjet.From, cfg.Mailjet.Name)

	examUC := usecase.NewExamUseCase(examRepo, questionRepo, taskRepo, geminiAdapter, publisher, cacheManager, mailer, taskCache, questionCache)
	questionUC := usecase.NewQuestionUseCase(questionRepo, questionCache, mailer)
	taskUC := usecase.NewTaskUseCase(taskRepo, taskCache)

	grpcServer, err := service.NewGRPCServer(*cfg, taskUC, questionUC, examUC)
	if err != nil {
		log.Printf("failed to initialize gRPC server: %v", err)
		return nil, fmt.Errorf("grpc server: %w", err)
	}
	log.Println("initialized gRPC server")

	return &App{
		grpcServer: grpcServer,
		mongo:      mongoClient,
		redis:      redisClient,
		nats:       natsConn,
		gemini:     geminiClient,
		tracer:     tracer,
	}, nil
}

func (a *App) Close() {
	log.Println("closing resources...")

	if a.grpcServer != nil {
		a.grpcServer.Stop()
	}

	if a.redis != nil {
		if err := a.redis.Close(); err != nil {
			log.Printf("failed to close Redis: %v", err)
		}
	}

	if a.nats != nil {
		a.nats.Close()
	}

	log.Println("all resources closed")
}

func (a *App) Run() error {
	errCh := make(chan error, 1)

	go func() {
		errCh <- a.grpcServer.Run()
	}()

	log.Printf("service %v started on port %d", serviceName, a.grpcServer.Cfg.Port)

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Printf("metrics server failed: %v", err)
			errCh <- err
		}
	}()
	log.Println("metrics server started on :9090")

	select {
	case err := <-errCh:
		log.Printf("grpc server failed: %v", err)
		return fmt.Errorf("grpc server failed: %w", err)
	case s := <-shutdownCh:
		log.Printf("received signal: %v. Running graceful shutdown...", s)
		a.Close()
		log.Println("graceful shutdown completed!")
	}

	return nil
}
