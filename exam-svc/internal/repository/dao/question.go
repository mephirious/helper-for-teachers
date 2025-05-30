package dao

import (
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Question struct {
	ID            primitive.ObjectID `bson:"_id"`
	ExamID        primitive.ObjectID `bson:"exam_id"`
	QuestionText  string             `bson:"question_text"`
	Options       []string           `bson:"options"`
	CorrectAnswer string             `bson:"correct_answer"`
	Status        string             `bson:"status"`
	CreatedAt     time.Time          `bson:"created_at"`
}

func FromDomainQuestion(q *domain.Question) *Question {
	return &Question{
		ID:            q.ID,
		ExamID:        q.ExamID,
		QuestionText:  q.QuestionText,
		Options:       q.Options,
		CorrectAnswer: q.CorrectAnswer,
		Status:        q.Status,
		CreatedAt:     q.CreatedAt,
	}
}

func (q *Question) ToDomainQuestion() *domain.Question {
	return &domain.Question{
		ID:            q.ID,
		ExamID:        q.ExamID,
		QuestionText:  q.QuestionText,
		Options:       q.Options,
		CorrectAnswer: q.CorrectAnswer,
		Status:        q.Status,
		CreatedAt:     q.CreatedAt,
	}
}
