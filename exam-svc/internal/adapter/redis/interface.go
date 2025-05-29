package redis

import (
	"context"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/repository"
)

type ExamCache interface {
	Init(ctx context.Context, repo repository.ExamRepository) error
	Set(ctx context.Context, exam domain.Exam) error
	Get(ctx context.Context, id string) (domain.Exam, error)
	Delete(ctx context.Context, id string) error
	SetMany(ctx context.Context, exams []domain.Exam) error
	GetAll(ctx context.Context) ([]domain.Exam, error)
}

type QuestionCache interface {
	Init(ctx context.Context, repo repository.QuestionRepository) error
	Set(ctx context.Context, question domain.Question) error
	Get(ctx context.Context, id string) (domain.Question, error)
	Delete(ctx context.Context, id string) error
	SetMany(ctx context.Context, questions []domain.Question) error
	GetAll(ctx context.Context) ([]domain.Question, error)
}

type TaskCache interface {
	Init(ctx context.Context, repo repository.TaskRepository) error
	Set(ctx context.Context, task domain.Task) error
	Get(ctx context.Context, id string) (domain.Task, error)
	Delete(ctx context.Context, id string) error
	SetMany(ctx context.Context, tasks []domain.Task) error
	GetAll(ctx context.Context) ([]domain.Task, error)
}
