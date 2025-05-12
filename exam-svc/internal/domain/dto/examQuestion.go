package dto

import (
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
)

type ExamQuestionCreate struct {
	ExamID        string   `json:"exam_id" bson:"exam_id"`
	QuestionText  string   `json:"question_text" bson:"question_text"`
	Options       []string `json:"options" bson:"options"`
	CorrectAnswer string   `json:"correct_answer" bson:"correct_answer"`
	Status        string   `json:"status" bson:"status"`
}

type ExamQuestionUpdate struct {
	ExamID        *string   `json:"exam_id" bson:"exam_id"`
	QuestionText  *string   `json:"question_text" bson:"question_text"`
	Options       []*string `json:"options" bson:"options"`
	CorrectAnswer *string   `json:"correct_answer" bson:"correct_answer"`
	Status        *string   `json:"status" bson:"status"`
}

type ExamQuestionPatch struct {
	Status *string `json:"status" bson:"status"`
}

type ExamQuestionResponse struct {
	ID            string    `json:"id" bson:"_id"`
	ExamID        string    `json:"exam_id" bson:"exam_id"`
	QuestionText  string    `json:"question_text" bson:"question_text"`
	Options       []string  `json:"options" bson:"options"`
	CorrectAnswer string    `json:"correct_answer" bson:"correct_answer"`
	Status        string    `json:"status" bson:"status"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
}

func MapExamQuestionToResponse(c domain.ExamQuestion) ExamQuestionResponse {
	return ExamQuestionResponse{
		ID:            c.ID.Hex(),
		ExamID:        c.ExamID.Hex(),
		QuestionText:  c.QuestionText,
		Options:       c.Options,
		CorrectAnswer: c.CorrectAnswer,
		Status:        c.Status,
		CreatedAt:     c.CreatedAt,
	}
}
