package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type questionUseCase struct {
	questionRepo repository.QuestionRepository
}

func NewQuestionUseCase(repo repository.QuestionRepository) QuestionUseCase {
	return &questionUseCase{
		questionRepo: repo,
	}
}

func (uc *questionUseCase) CreateQuestion(ctx context.Context, question *domain.Question) (*domain.Question, error) {
	question.ID = primitive.NewObjectID()
	question.CreatedAt = time.Now()

	if err := uc.questionRepo.CreateQuestion(ctx, question); err != nil {
		return nil, fmt.Errorf("failed to create question: %w", err)
	}
	return question, nil
}

func (uc *questionUseCase) GetQuestionByID(ctx context.Context, id primitive.ObjectID) (*domain.Question, error) {
	q, err := uc.questionRepo.GetQuestionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if q == nil {
		return nil, fmt.Errorf("question not found")
	}
	return q, nil
}

func (uc *questionUseCase) GetQuestionsByExamID(ctx context.Context, examID primitive.ObjectID) ([]domain.Question, error) {
	return uc.questionRepo.GetQuestionsByExamID(ctx, examID)
}

func (uc *questionUseCase) GetAllQuestions(ctx context.Context) ([]domain.Question, error) {
	return uc.questionRepo.GetAllQuestions(ctx)
}

func (uc *questionUseCase) DeleteQuestion(ctx context.Context, id primitive.ObjectID) error {
	return uc.questionRepo.DeleteQuestion(ctx, id)
}

func (u *questionUseCase) UpdateQuestion(ctx context.Context, question *domain.Question) error {
	return u.questionRepo.UpdateQuestion(ctx, question)
}
