package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type examUseCase struct {
	examRepo     repository.ExamRepository
	questionRepo repository.QuestionRepository
	taskRepo     repository.TaskRepository
}

func NewExamUseCase(examRepo repository.ExamRepository, questionRepo repository.QuestionRepository, taskRepo repository.TaskRepository) ExamUseCase {
	return &examUseCase{
		examRepo:     examRepo,
		questionRepo: questionRepo,
		taskRepo:     taskRepo,
	}
}

func (uc *examUseCase) CreateExam(ctx context.Context, exam *domain.Exam) (*domain.Exam, error) {
	exam.ID = primitive.NewObjectID()
	exam.CreatedAt = time.Now()
	exam.UpdatedAt = exam.CreatedAt

	if err := uc.examRepo.CreateExam(ctx, exam); err != nil {
		return nil, fmt.Errorf("failed to create exam: %w", err)
	}
	return exam, nil
}

func (uc *examUseCase) GetExamByID(ctx context.Context, id primitive.ObjectID) (*domain.Exam, error) {
	exam, err := uc.examRepo.GetExamByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if exam == nil {
		return nil, fmt.Errorf("exam not found")
	}
	return exam, nil
}

func (uc *examUseCase) GetExamsByUser(ctx context.Context, userID primitive.ObjectID) ([]domain.Exam, error) {
	return uc.examRepo.GetExamsByUser(ctx, userID)
}

func (uc *examUseCase) UpdateExamStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	return uc.examRepo.UpdateExamStatus(ctx, id, status)
}

func (uc *examUseCase) DeleteExam(ctx context.Context, id primitive.ObjectID) error {
	return uc.examRepo.DeleteExam(ctx, id)
}

func (uc *examUseCase) GetAllExams(ctx context.Context) ([]domain.Exam, error) {
	return uc.examRepo.GetAllExams(ctx)
}

func (uc *examUseCase) GetExamWithDetails(ctx context.Context, id primitive.ObjectID) (*domain.ExamDetailed, error) {
	exam, err := uc.examRepo.GetExamByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if exam == nil {
		return nil, fmt.Errorf("exam not found")
	}

	questionsResult, err := uc.questionRepo.GetQuestionsByExamID(ctx, id)
	if err != nil {
		return nil, err
	}

	tasksResult, err := uc.taskRepo.GetTasksByExamID(ctx, id)
	if err != nil {
		return nil, err
	}

	var questions []domain.Question
	for _, q := range questionsResult {
		questions = append(questions, domain.Question{
			ID:            q.ID,
			ExamID:        q.ExamID,
			QuestionText:  q.QuestionText,
			Options:       q.Options,
			CorrectAnswer: q.CorrectAnswer,
			Status:        q.Status,
			CreatedAt:     q.CreatedAt,
		})
	}

	var tasks []domain.Task
	for _, t := range tasksResult {
		tasks = append(tasks, domain.Task{
			ID:          t.ID,
			ExamID:      t.ExamID,
			TaskType:    t.TaskType,
			Description: t.Description,
			Score:       t.Score,
			CreatedAt:   t.CreatedAt,
		})
	}

	return &domain.ExamDetailed{
		ID:          exam.ID,
		Title:       exam.Title,
		Description: exam.Description,
		CreatedBy:   exam.CreatedBy,
		Status:      exam.Status,
		CreatedAt:   exam.CreatedAt,
		UpdatedAt:   exam.UpdatedAt,
		Tasks:       tasks,
		Questions:   questions,
	}, nil
}
