package dto

import (
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	pb "github.com/mephirious/helper-for-teachers/services/exam-svc/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ExamCreate struct {
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	CreatedBy   string `json:"created_by" bson:"created_by"`
	Status      string `json:"status" bson:"status"`
}

type ExamCreateAI struct {
	Grade string `json:"grade" bson:"grade"`
	Topic string `json:"topic" bson:"topic"`
}

type ExamUpdate struct {
	Title       *string `json:"title" bson:"title"`
	Description *string `json:"description" bson:"description"`
	CreatedBy   *string `json:"created_by" bson:"created_by"`
	Status      *string `json:"status" bson:"status"`
}

type ExamPatch struct {
	Status *string `json:"status" bson:"status"`
}

type ExamResponse struct {
	ID          string    `json:"id" bson:"_id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	CreatedBy   string    `json:"created_by" bson:"created_by"`
	Status      string    `json:"status" bson:"status"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

type ExamDetailedResponse struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	CreatedBy   string             `json:"created_by"`
	Status      string             `json:"status"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	Tasks       []TaskResponse     `json:"tasks"`
	Questions   []QuestionResponse `json:"questions"`
}

func MapExamToResponse(c domain.Exam) ExamResponse {
	return ExamResponse{
		ID:          c.ID.Hex(),
		Title:       c.Title,
		Description: c.Description,
		CreatedBy:   c.CreatedBy.Hex(),
		Status:      c.Status,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

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
