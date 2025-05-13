package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/gemini"
	nats "github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/nats"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/repository"
	pb "github.com/mephirious/helper-for-teachers/services/exam-svc/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type examUseCase struct {
	examRepo     repository.ExamRepository
	questionRepo repository.QuestionRepository
	taskRepo     repository.TaskRepository
	geminiClient *gemini.Client
	publisher    nats.ExamEventProducer
}

func NewExamUseCase(examRepo repository.ExamRepository, questionRepo repository.QuestionRepository, taskRepo repository.TaskRepository, geminiClient *gemini.Client) ExamUseCase {
	return &examUseCase{
		examRepo:     examRepo,
		questionRepo: questionRepo,
		taskRepo:     taskRepo,
		geminiClient: geminiClient,
	}
}

func (uc *examUseCase) CreateExam(ctx context.Context, exam *domain.Exam) (*domain.Exam, error) {
	exam.ID = primitive.NewObjectID()
	exam.CreatedAt = time.Now()
	exam.UpdatedAt = exam.CreatedAt

	if err := uc.examRepo.CreateExam(ctx, exam); err != nil {
		return nil, fmt.Errorf("failed to create exam: %w", err)
	}

	if err := uc.publisher.Push(ctx, exam, pb.ExamEventType_CREATED); err != nil {
		log.Printf("Failed to push create event to NATS: %v", err)
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
	exam, err := uc.examRepo.GetExamByID(ctx, id)
	if err != nil {
		return err
	}
	if exam == nil {
		return fmt.Errorf("exam not found")
	}

	exam.Status = status
	exam.UpdatedAt = time.Now()

	if err := uc.examRepo.UpdateExamStatus(ctx, id, status); err != nil {
		return err
	}

	if err := uc.publisher.Push(ctx, exam, pb.ExamEventType_UPDATED); err != nil {
		log.Printf("Failed to push update event to NATS: %v", err)
	}

	return nil
}

func (uc *examUseCase) UpdateExam(ctx context.Context, exam *domain.Exam) error {
	exam.UpdatedAt = time.Now()

	if err := uc.examRepo.UpdateExam(ctx, exam); err != nil {
		return err
	}

	updated, err := uc.examRepo.GetExamByID(ctx, exam.ID)
	if err != nil || updated == nil {
		log.Printf("Failed to reload exam for NATS publish: %v", err)
		return nil
	}

	if err := uc.publisher.Push(ctx, updated, pb.ExamEventType_UPDATED); err != nil {
		log.Printf("Failed to push update event to NATS: %v", err)
	}

	return nil
}

func (uc *examUseCase) DeleteExam(ctx context.Context, id primitive.ObjectID) error {
	exam, err := uc.examRepo.GetExamByID(ctx, id)
	if err != nil {
		return err
	}
	if exam == nil {
		return fmt.Errorf("exam not found")
	}

	if err := uc.examRepo.DeleteExam(ctx, id); err != nil {
		return err
	}

	if err := uc.publisher.Push(ctx, exam, pb.ExamEventType_DELETED); err != nil {
		log.Printf("Failed to push delete event to NATS: %v", err)
	}

	return nil
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

func (uc *examUseCase) GenerateExamUsingAI(ctx context.Context, userID primitive.ObjectID, numQuestions, numTasks int, topic, grade string) (*domain.ExamDetailed, error) {
	result, err := uc.geminiClient.GenerateExam(ctx, numQuestions, numTasks, grade, topic)
	if err != nil {
		return nil, fmt.Errorf("failed to generate exam using AI: %w", err)
	}

	examID := primitive.NewObjectID()
	now := time.Now()

	exam := &domain.Exam{
		ID:          examID,
		Title:       result.Title,
		Description: result.Description,
		Status:      result.Status,
		CreatedBy:   userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.examRepo.CreateExam(ctx, exam); err != nil {
		return nil, fmt.Errorf("failed to save generated exam: %w", err)
	}

	var questions []domain.Question
	for _, q := range result.Questions {
		q.ID = primitive.NewObjectID()
		q.ExamID = examID
		q.CreatedAt = now

		if err := uc.questionRepo.CreateQuestion(ctx, &q); err != nil {
			return nil, fmt.Errorf("failed to save question: %w", err)
		}
		questions = append(questions, q)
	}

	var tasks []domain.Task
	for _, t := range result.Tasks {
		t.ID = primitive.NewObjectID()
		t.ExamID = examID
		t.CreatedAt = now

		if err := uc.taskRepo.CreateTask(ctx, &t); err != nil {
			return nil, fmt.Errorf("failed to save task: %w", err)
		}
		tasks = append(tasks, t)
	}

	return &domain.ExamDetailed{
		ID:          exam.ID,
		Title:       exam.Title,
		Description: exam.Description,
		CreatedBy:   exam.CreatedBy,
		Status:      exam.Status,
		CreatedAt:   exam.CreatedAt,
		UpdatedAt:   exam.UpdatedAt,
		Questions:   questions,
		Tasks:       tasks,
	}, nil
}
