package nats

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/mephirious/helper-for-teachers/managers-svc/internal/core/domain"
	group_eventpb "github.com/mephirious/helper-for-teachers/managers-svc/internal/grpc"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type EventPublisher struct {
	conn    *nats.Conn
	timeout time.Duration
	logger  *slog.Logger
}

func NewEventPublisher(natsURL string, logger *slog.Logger) (*EventPublisher, error) {
	nc, err := nats.Connect(natsURL, nats.ErrorHandler(func(nc *nats.Conn, s *nats.Subscription, err error) {
		logger.Error("NATS Error", "error", err.Error())
	}))
	if err != nil {
		return nil, fmt.Errorf("NATS connection failed: %w", err)
	}
	logger.Info("Connected to NATS", "at", nc.ConnectedUrl())

	return &EventPublisher{
		conn:    nc,
		timeout: 2 * time.Second,
		logger:  logger,
	}, nil
}

func (p *EventPublisher) PublishCourseCreated(ctx context.Context, course *domain.Course) error {
	event := &group_eventpb.CourseCreatedEvent{
		Id:        course.ID.String(),
		Name:      course.Name,
		CreatedAt: course.CreatedAt.Unix(),
		UpdatedAt: course.UpdatedAt.Unix(),
	}
	return p.publishEvent(ctx, "course.created", event)
}

func (p *EventPublisher) PublishGroupCreated(ctx context.Context, group *domain.Group) error {
	event := &group_eventpb.GroupCreatedEvent{
		Id:        group.ID.String(),
		CourseId:  group.CourseID.String(),
		Name:      group.Name,
		ExpireAt:  group.ExpireAt.Unix(),
		CreatedAt: group.CreatedAt.Unix(),
		UpdatedAt: group.UpdatedAt.Unix(),
	}
	return p.publishEvent(ctx, "group.created", event)
}

func (p *EventPublisher) PublishGroupMemberAdded(ctx context.Context, member *domain.GroupMember) error {
	event := &group_eventpb.GroupMemberAddedEvent{
		Id:        member.ID.String(),
		GroupId:   member.GroupID.String(),
		UserId:    member.UserID.String(),
		Role:      string(member.Role),
		CreatedAt: member.CreatedAt.Unix(),
		UpdatedAt: member.UpdatedAt.Unix(),
	}
	return p.publishEvent(ctx, "group.member.added", event)
}

func (p *EventPublisher) PublishCourseInstructorAssigned(ctx context.Context, instr *domain.CourseInstructor) error {
	event := &group_eventpb.CourseInstructorAssignedEvent{
		Id:        instr.ID.String(),
		CourseId:  instr.CourseID.String(),
		UserId:    instr.UserID.String(),
		CreatedAt: instr.CreatedAt.Unix(),
		UpdatedAt: instr.UpdatedAt.Unix(),
	}
	return p.publishEvent(ctx, "course.instructor.assigned", event)
}

func (p *EventPublisher) publishEvent(ctx context.Context, subject string, event proto.Message) error {
	data, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	if err := p.conn.Publish(subject, data); err != nil {
		p.logger.Error("Failed to publish", "to", subject, "error", err.Error())
		return fmt.Errorf("publish failed: %w", err)
	}
	p.logger.Info("Event published", "to", subject, "event", event)
	return nil
}

func (p *EventPublisher) Close() {
	p.conn.Close()
}
