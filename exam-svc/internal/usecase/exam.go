package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/gemini"
	inmemory "github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/in-memory"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/mailjet"
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
	publisher    *nats.ExamEventProducer
	cache        *inmemory.CacheManager
	mailjet      *mailjet.MailjetClient
}

func NewExamUseCase(examRepo repository.ExamRepository, questionRepo repository.QuestionRepository, taskRepo repository.TaskRepository, geminiClient *gemini.Client, publisher *nats.ExamEventProducer, cache *inmemory.CacheManager, mailjetClient *mailjet.MailjetClient) ExamUseCase {
	return &examUseCase{
		examRepo:     examRepo,
		questionRepo: questionRepo,
		taskRepo:     taskRepo,
		geminiClient: geminiClient,
		publisher:    publisher,
		cache:        cache,
		mailjet:      mailjetClient,
	}
}

func (uc *examUseCase) CreateExam(ctx context.Context, exam *domain.Exam) (*domain.Exam, error) {
	exam.ID = primitive.NewObjectID()
	exam.CreatedAt = time.Now()
	exam.UpdatedAt = exam.CreatedAt

	if err := uc.examRepo.CreateExam(ctx, exam); err != nil {
		return nil, fmt.Errorf("failed to create exam: %w", err)
	}

	uc.cache.ExamCache.Set(*exam)

	if err := uc.publisher.Push(ctx, exam, pb.ExamEventType_EXAM_CREATED); err != nil {
		log.Printf("Failed to push create event to NATS: %v", err)
	}

	return exam, nil
}

func (uc *examUseCase) GetExamByID(ctx context.Context, id primitive.ObjectID) (*domain.Exam, error) {
	exam, ok := uc.cache.ExamCache.Get(id.Hex())
	if ok {
		return &exam, nil
	}

	examPtr, err := uc.examRepo.GetExamByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if examPtr == nil {
		return nil, fmt.Errorf("exam not found")
	}

	uc.cache.ExamCache.Set(*examPtr)
	return examPtr, nil
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

	uc.cache.ExamCache.Set(*exam)

	if err := uc.publisher.Push(ctx, exam, pb.ExamEventType_EXAM_UPDATED); err != nil {
		log.Printf("Failed to push update event to NATS: %v", err)
	}

	if status == "verified" && uc.mailjet != nil {
		if err := uc.mailjet.SendTemplateEmail("admin@example.com", "Admin", mailjet.ExamVerifiedTemplate); err != nil {
			log.Printf("Failed to send exam verified email: %v", err)
		} else {
			log.Printf("Sent exam verified email for exam ID: %s", id.Hex())
		}
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

	uc.cache.ExamCache.Set(*updated)

	if err := uc.publisher.Push(ctx, updated, pb.ExamEventType_EXAM_UPDATED); err != nil {
		log.Printf("Failed to push update event to NATS: %v", err)
	}

	if exam.Status == "verified" && uc.mailjet != nil {
		if err := uc.mailjet.SendTemplateEmail("admin@example.com", "Admin", mailjet.ExamVerifiedTemplate); err != nil {
			log.Printf("Failed to send exam verified email: %v", err)
		} else {
			log.Printf("Sent exam verified email for exam ID: %s", exam.ID.Hex())
		}
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

	uc.cache.ExamCache.Delete(id.Hex())

	if err := uc.publisher.Push(ctx, exam, pb.ExamEventType_EXAM_DELETED); err != nil {
		log.Printf("Failed to push delete event to NATS: %v", err)
	}

	return nil
}

func (uc *examUseCase) GetAllExams(ctx context.Context) ([]domain.Exam, error) {
	exams := uc.cache.ExamCache.GetAll()
	if len(exams) > 0 {
		return exams, nil
	}

	exams, err := uc.examRepo.GetAllExams(ctx)
	if err != nil {
		return nil, err
	}

	uc.cache.ExamCache.SetMany(exams)
	return exams, nil
}

func (uc *examUseCase) GetExamWithDetails(ctx context.Context, id primitive.ObjectID) (*domain.ExamDetailed, error) {
	exam, err := uc.GetExamByID(ctx, id)
	if err != nil {
		return nil, err
	}

	questions, err := uc.questionRepo.GetQuestionsByExamID(ctx, id)
	if err != nil {
		return nil, err
	}

	tasks, err := uc.taskRepo.GetTasksByExamID(ctx, id)
	if err != nil {
		return nil, err
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

	uc.cache.ExamCache.Set(*exam)

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
