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

type QuestionUseCase interface {
	CreateQuestion(ctx context.Context, question *domain.Question) (*domain.Question, error)
	GetQuestionByID(ctx context.Context, id primitive.ObjectID) (*domain.Question, error)
	GetQuestionsByExamID(ctx context.Context, examID primitive.ObjectID) ([]domain.Question, error)
	GetAllQuestions(ctx context.Context) ([]domain.Question, error)
	DeleteQuestion(ctx context.Context, id primitive.ObjectID) error
}

type ExamUseCase interface {
	CreateExam(ctx context.Context, exam *domain.Exam) (*domain.Exam, error)
	GetExamByID(ctx context.Context, id primitive.ObjectID) (*domain.Exam, error)
	GetExamsByUser(ctx context.Context, userID primitive.ObjectID) ([]domain.Exam, error)
	UpdateExamStatus(ctx context.Context, id primitive.ObjectID, status string) error
	DeleteExam(ctx context.Context, id primitive.ObjectID) error
	GetAllExams(ctx context.Context) ([]domain.Exam, error)
	GetExamWithDetails(ctx context.Context, id primitive.ObjectID) (*domain.ExamDetailed, error)
	GenerateExamUsingAI(ctx context.Context, userID primitive.ObjectID, numQuestions, numTasks int, topic, grade string) (*domain.ExamDetailed, error)
}
