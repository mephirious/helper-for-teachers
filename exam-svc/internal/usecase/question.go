package usecase

import (
	"context"
	"fmt"
	"time"

	redis "github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/redis/cache"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type questionUseCase struct {
	questionRepo  repository.QuestionRepository
	questionCache redis.QuestionCache
}

func NewQuestionUseCase(repo repository.QuestionRepository, cache redis.QuestionCache) QuestionUseCase {
	return &questionUseCase{
		questionRepo:  repo,
		questionCache: cache,
	}
}

func (uc *questionUseCase) CreateQuestion(ctx context.Context, question *domain.Question) (*domain.Question, error) {
	question.ID = primitive.NewObjectID()
	question.CreatedAt = time.Now()

	if err := uc.questionRepo.CreateQuestion(ctx, question); err != nil {
		return nil, fmt.Errorf("failed to create question: %w", err)
	}

	_ = uc.questionCache.Set(ctx, *question)
	return question, nil
}

func (uc *questionUseCase) GetQuestionByID(ctx context.Context, id primitive.ObjectID) (*domain.Question, error) {
	q, err := uc.questionCache.Get(ctx, id.Hex())
	if err == nil {
		return &q, nil
	}

	qptr, err := uc.questionRepo.GetQuestionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if qptr == nil {
		return nil, fmt.Errorf("question not found")
	}

	_ = uc.questionCache.Set(ctx, *qptr)
	return qptr, nil
}

func (uc *questionUseCase) GetQuestionsByExamID(ctx context.Context, examID primitive.ObjectID) ([]domain.Question, error) {
	return uc.questionRepo.GetQuestionsByExamID(ctx, examID)
}

func (uc *questionUseCase) GetAllQuestions(ctx context.Context) ([]domain.Question, error) {
	questions, err := uc.questionCache.GetAll(ctx)
	if err == nil && len(questions) > 0 {
		return questions, nil
	}

	questions, err = uc.questionRepo.GetAllQuestions(ctx)
	if err != nil {
		return nil, err
	}

	_ = uc.questionCache.SetMany(ctx, questions)
	return questions, nil
}

func (uc *questionUseCase) DeleteQuestion(ctx context.Context, id primitive.ObjectID) error {
	if err := uc.questionRepo.DeleteQuestion(ctx, id); err != nil {
		return err
	}
	_ = uc.questionCache.Delete(ctx, id.Hex())
	return nil
}

func (uc *questionUseCase) UpdateQuestion(ctx context.Context, question *domain.Question) error {
	if err := uc.questionRepo.UpdateQuestion(ctx, question); err != nil {
		return err
	}
	_ = uc.questionCache.Set(ctx, *question)
	return nil
}
