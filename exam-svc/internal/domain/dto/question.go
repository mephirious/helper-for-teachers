package dto

import (
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
)

type QuestionCreate struct {
	ExamID        string   `json:"exam_id" bson:"exam_id"`
	QuestionText  string   `json:"question_text" bson:"question_text"`
	Options       []string `json:"options" bson:"options"`
	CorrectAnswer string   `json:"correct_answer" bson:"correct_answer"`
	Status        string   `json:"status" bson:"status"`
}

type QuestionUpdate struct {
	ExamID        *string   `json:"exam_id" bson:"exam_id"`
	QuestionText  *string   `json:"question_text" bson:"question_text"`
	Options       []*string `json:"options" bson:"options"`
	CorrectAnswer *string   `json:"correct_answer" bson:"correct_answer"`
	Status        *string   `json:"status" bson:"status"`
}

type QuestionPatch struct {
	Status *string `json:"status" bson:"status"`
}

type QuestionResponse struct {
	ID            string    `json:"id" bson:"_id"`
	ExamID        string    `json:"exam_id" bson:"exam_id"`
	QuestionText  string    `json:"question_text" bson:"question_text"`
	Options       []string  `json:"options" bson:"options"`
	CorrectAnswer string    `json:"correct_answer" bson:"correct_answer"`
	Status        string    `json:"status" bson:"status"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
}

func MapQuestionToResponse(c domain.Question) QuestionResponse {
	return QuestionResponse{
		ID:            c.ID.Hex(),
		ExamID:        c.ExamID.Hex(),
		QuestionText:  c.QuestionText,
		Options:       c.Options,
		CorrectAnswer: c.CorrectAnswer,
		Status:        c.Status,
		CreatedAt:     c.CreatedAt,
	}
}
