package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	eventspb "github.com/suyundykovv/margulan-protos/gen/go/events/v1"
	"go.uber.org/zap"
)

const (
	connectTimeout = 5 * time.Second
	lessonSubject  = "school.lessons.%s"
	taskSubject    = "school.tasks.%s"
)

type SchoolEventPublisher struct {
	conn   *nats.Conn
	js     nats.JetStreamContext
	logger *zap.Logger
}

func NewSchoolEventPublisher(natsURL string, logger *zap.Logger) (*SchoolEventPublisher, error) {
	opts := []nats.Option{
		nats.Timeout(connectTimeout),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(5),
	}

	conn, err := nats.Connect(natsURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	streamConfigs := []*nats.StreamConfig{
		{
			Name:     "SCHOOL_LESSONS",
			Subjects: []string{"school.lessons.>"},
			MaxAge:   24 * time.Hour,
		},
		{
			Name:     "SCHOOL_TASKS",
			Subjects: []string{"school.tasks.>"},
			MaxAge:   72 * time.Hour,
		},
	}

	for _, cfg := range streamConfigs {
		_, err = js.AddStream(cfg)
		if err != nil {
			logger.Warn("Stream creation failed (may already exist)",
				zap.String("stream", cfg.Name),
				zap.Error(err))
		}
	}

	return &SchoolEventPublisher{
		conn:   conn,
		js:     js,
		logger: logger,
	}, nil
}

func (p *SchoolEventPublisher) PublishLessonCreated(ctx context.Context, lesson *eventspb.Lesson) error {
	return p.publishLesson(ctx, "created", lesson)
}

func (p *SchoolEventPublisher) PublishLessonUpdated(ctx context.Context, lesson *eventspb.Lesson) error {
	return p.publishLesson(ctx, "updated", lesson)
}

func (p *SchoolEventPublisher) PublishLessonDeleted(ctx context.Context, lessonID string) error {
	return p.publishLesson(ctx, "deleted", map[string]string{"lesson_id": lessonID})
}

func (p *SchoolEventPublisher) PublishTaskCreated(ctx context.Context, task *eventspb.Task) error {
	return p.publishTask(ctx, "created", task)
}

func (p *SchoolEventPublisher) PublishTaskUpdated(ctx context.Context, task *eventspb.Task) error {
	return p.publishTask(ctx, "updated", task)
}

func (p *SchoolEventPublisher) PublishTaskDeleted(ctx context.Context, taskID string) error {
	return p.publishTask(ctx, "deleted", map[string]string{"task_id": taskID})
}

func (p *SchoolEventPublisher) publishLesson(ctx context.Context, action string, data interface{}) error {
	return p.publish(ctx, fmt.Sprintf(lessonSubject, action), data)
}

func (p *SchoolEventPublisher) publishTask(ctx context.Context, action string, data interface{}) error {
	return p.publish(ctx, fmt.Sprintf(taskSubject, action), data)
}

func (p *SchoolEventPublisher) publish(ctx context.Context, subject string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	pubAck, err := p.js.Publish(subject, jsonData)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	p.logger.Debug("Message published",
		zap.String("subject", subject),
		zap.Uint64("sequence", pubAck.Sequence))

	return nil
}

func (p *SchoolEventPublisher) Close() {
	if p.conn != nil && !p.conn.IsClosed() {
		p.conn.Close()
	}
}
