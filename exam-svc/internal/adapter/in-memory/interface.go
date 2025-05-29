package inmemory

import (
	"context"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/repository"
)

type ExamCacheInterface interface {
	Init(ctx context.Context, repo repository.ExamRepository) error
	Set(domain.Exam)
	SetMany([]domain.Exam)
	Get(string) (domain.Exam, bool)
	GetAll() []domain.Exam
	Delete(string)
}

type QuestionCacheInterface interface {
	Init(ctx context.Context, repo repository.QuestionRepository) error
	Set(domain.Question)
	SetMany([]domain.Question)
	Get(string) (domain.Question, bool)
	GetAll() []domain.Question
	Delete(string)
}

type TaskCacheInterface interface {
	Init(ctx context.Context, repo repository.TaskRepository) error
	Set(domain.Task)
	SetMany([]domain.Task)
	Get(string) (domain.Task, bool)
	GetAll() []domain.Task
	Delete(string)
}
