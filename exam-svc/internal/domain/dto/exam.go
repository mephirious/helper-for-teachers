package dto

import (
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
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
