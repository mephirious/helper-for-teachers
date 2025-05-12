package inmemory

import "github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"

type ExamCacheInterface interface {
	Set(domain.Exam)
	SetMany([]domain.Exam)
	Get(string) (domain.Exam, bool)
	GetAll() []domain.Exam
	Delete(string)
}

type QuestionCacheInterface interface {
	Set(domain.Question)
	SetMany([]domain.Question)
	Get(string) (domain.Question, bool)
	GetAll() []domain.Question
	Delete(string)
}

type TaskCacheInterface interface {
	Set(domain.Task)
	SetMany([]domain.Task)
	Get(string) (domain.Task, bool)
	GetAll() []domain.Task
	Delete(string)
}
