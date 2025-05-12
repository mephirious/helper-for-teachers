package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	ExamID      primitive.ObjectID `json:"exam_id" bson:"exam_id"`
	TaskType    string             `json:"task_type" bson:"task_type"`
	Description string             `json:"description" bson:"description"`
	Score       float32            `json:"score" bson:"score"`
	CreatedAt   time.Time          `json:"created_At" bson:"created_At"`
}
