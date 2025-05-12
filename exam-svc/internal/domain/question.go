package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Question struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	ExamID        primitive.ObjectID `json:"exam_id" bson:"exam_id"`
	QuestionText  string             `json:"question_text" bson:"question_text"`
	Options       []string           `json:"options" bson:"options"`
	CorrectAnswer string             `json:"correct_answer" bson:"correct_answer"`
	Status        string             `json:"status" bson:"status"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
}
