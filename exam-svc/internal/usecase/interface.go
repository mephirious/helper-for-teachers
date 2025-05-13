package usecase

import (
	"context"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskUseCase interface {
	CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error)
	GetTaskByID(ctx context.Context, id primitive.ObjectID) (*domain.Task, error)
	GetTasksByExamID(ctx context.Context, examID primitive.ObjectID) ([]domain.Task, error)
	GetAllTasks(ctx context.Context) ([]domain.Task, error)
	DeleteTask(ctx context.Context, id primitive.ObjectID) error
}
