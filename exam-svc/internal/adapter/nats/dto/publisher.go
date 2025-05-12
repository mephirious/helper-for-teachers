package dto

import (
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	pb "github.com/mephirious/helper-for-teachers/services/exam-svc/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToExamEvent(exam *domain.Exam, eventType pb.ExamEventType) *pb.ExamEvent {
	return &pb.ExamEvent{
		Id:          exam.ID.Hex(),
		Title:       exam.Title,
		Description: exam.Description,
		CreatedBy:   exam.CreatedBy.Hex(),
		Status:      exam.Status,
		CreatedAt:   timestamppb.New(exam.CreatedAt),
		UpdatedAt:   timestamppb.New(exam.UpdatedAt),
		EventType:   eventType,
	}
}
