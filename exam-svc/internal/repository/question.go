package repository

import (
	"context"
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/repository/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type questionRepository struct {
	collection *mongo.Collection
}

func NewQuestionRepository(db *mongo.Database) QuestionRepository {
	return &questionRepository{
		collection: db.Collection("questions"),
	}
}

func (r *questionRepository) CreateQuestion(ctx context.Context, question *domain.Question) error {
	question.ID = primitive.NewObjectID()
	question.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, dao.FromDomainQuestion(question))
	return err
}

func (r *questionRepository) GetQuestionByID(ctx context.Context, id primitive.ObjectID) (*domain.Question, error) {
	var questionDAO dao.Question
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&questionDAO)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return questionDAO.ToDomainQuestion(), nil
}

func (r *questionRepository) GetQuestionsByExamID(ctx context.Context, examID primitive.ObjectID) ([]domain.Question, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"exam_id": examID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var questionDAOs []dao.Question
	if err := cursor.All(ctx, &questionDAOs); err != nil {
		return nil, err
	}

	questions := make([]domain.Question, 0, len(questionDAOs))
	for _, q := range questionDAOs {
		questions = append(questions, *q.ToDomainQuestion())
	}
	return questions, nil
}

func (r *questionRepository) GetAllQuestions(ctx context.Context) ([]domain.Question, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var questionDAOs []dao.Question
	if err := cursor.All(ctx, &questionDAOs); err != nil {
		return nil, err
	}

	questions := make([]domain.Question, 0, len(questionDAOs))
	for _, q := range questionDAOs {
		questions = append(questions, *q.ToDomainQuestion())
	}
	return questions, nil
}

func (r *questionRepository) DeleteQuestion(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
