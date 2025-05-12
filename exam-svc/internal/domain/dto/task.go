package dto

import (
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
)

type TaskCreate struct {
	ExamID      string  `json:"exam_id" bson:"exam_id"`
	TaskType    string  `json:"task_type" bson:"task_type"`
	Description string  `json:"description" bson:"description"`
	Score       float32 `json:"score" bson:"score"`
}

type TaskUpdate struct {
	ExamID      *string  `json:"exam_id" bson:"exam_id"`
	TaskType    *string  `json:"task_type" bson:"task_type"`
	Description *string  `json:"description" bson:"description"`
	Score       *float32 `json:"score" bson:"score"`
}

type TaskPatch struct {
	Score *float32 `json:"score" bson:"score"`
}

type TaskResponse struct {
	ID          string    `json:"id" bson:"_id"`
	ExamID      string    `json:"exam_id" bson:"exam_id"`
	TaskType    string    `json:"task_type" bson:"task_type"`
	Description string    `json:"description" bson:"description"`
	Score       float32   `json:"score" bson:"score"`
	CreatedAt   time.Time `json:"created_At" bson:"created_At"`
}

func MapTaskToResponse(c domain.Task) TaskResponse {
	return TaskResponse{
		ID:          c.ID.Hex(),
		ExamID:      c.ExamID.Hex(),
		TaskType:    c.TaskType,
		Description: c.Description,
		Score:       c.Score,
		CreatedAt:   c.CreatedAt,
	}
}
