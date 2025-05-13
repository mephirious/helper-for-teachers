package dao

import (
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id"`
	ExamID      primitive.ObjectID `bson:"exam_id"`
	TaskType    string             `bson:"task_type"`
	Description string             `bson:"description"`
	Score       float32            `bson:"score"`
	CreatedAt   time.Time          `bson:"created_at"`
}

func (t *Task) ToDomain() *domain.Task {
	return &domain.Task{
		ID:          t.ID,
		ExamID:      t.ExamID,
		TaskType:    t.TaskType,
		Description: t.Description,
		Score:       t.Score,
		CreatedAt:   t.CreatedAt,
	}
}

func FromDomain(task *domain.Task) *Task {
	return &Task{
		ID:          task.ID,
		ExamID:      task.ExamID,
		TaskType:    task.TaskType,
		Description: task.Description,
		Score:       task.Score,
		CreatedAt:   task.CreatedAt,
	}
}
