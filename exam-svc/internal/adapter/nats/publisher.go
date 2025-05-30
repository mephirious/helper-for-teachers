package producer

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/pkg/nats"
	pb "github.com/mephirious/helper-for-teachers/services/exam-svc/proto"
)

const PushTimeout = time.Second * 30

type ExamEventProducer struct {
	natsClient *nats.Client
	subject    string
}

func NewExamEventProducer(natsClient *nats.Client, subject string) *ExamEventProducer {
	return &ExamEventProducer{
		natsClient: natsClient,
		subject:    subject,
	}
}

func (p *ExamEventProducer) Push(ctx context.Context, exam *domain.Exam, eventType pb.ExamEventType) error {
	ctx, cancel := context.WithTimeout(ctx, PushTimeout)
	defer cancel()

	if eventType == pb.ExamEventType_EXAM_CREATED && exam.Title == "" {
		return fmt.Errorf("missing title for CREATED exam event")
	}

	pbEvent := &pb.ExamEvent{
		Id:          exam.ID.Hex(),
		Title:       exam.Title,
		Description: exam.Description,
		CreatedBy:   exam.CreatedBy.Hex(),
		Status:      exam.Status,
		CreatedAt:   timestamppb.New(exam.CreatedAt),
		UpdatedAt:   timestamppb.New(exam.UpdatedAt),
		EventType:   eventType,
	}

	data, err := proto.Marshal(pbEvent)
	if err != nil {
		return fmt.Errorf("proto.Marshal: %w", err)
	}

	log.Printf("Publishing ExamEvent to %s: %+v, data: %s", p.subject, pbEvent, hex.EncodeToString(data))
	err = p.natsClient.Conn.Publish(p.subject, data)
	if err != nil {
		return fmt.Errorf("p.natsClient.Conn.Publish: %w", err)
	}
	log.Printf("Exam event pushed to %s: %+v [%s]", p.subject, exam, eventType)

	return nil
}
