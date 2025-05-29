package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/repository/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type questionRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

func NewQuestionRepository(db *mongo.Database, client *mongo.Client) QuestionRepository {
	return &questionRepository{
		collection: db.Collection("questions"),
		client:     client,
	}
}

func (r *questionRepository) CreateQuestion(ctx context.Context, question *domain.Question) error {
	session, err := r.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
		examColl := r.collection.Database().Collection("exams")
		var examDAO dao.Exam
		err := examColl.FindOne(sessionCtx, bson.M{"_id": question.ExamID}).Decode(&examDAO)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("exam with id %s not found", question.ExamID.Hex())
			}
			return nil, fmt.Errorf("failed to verify exam: %w", err)
		}

		question.ID = primitive.NewObjectID()
		question.CreatedAt = time.Now()
		_, err = r.collection.InsertOne(sessionCtx, dao.FromDomainQuestion(question))
		if err != nil {
			return nil, fmt.Errorf("failed to insert question: %w", err)
		}

		return nil, nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
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

func (r *questionRepository) UpdateQuestion(ctx context.Context, question *domain.Question) error {
	update := bson.M{}
	if question.QuestionText != "" {
		update["question_text"] = question.QuestionText
	}
	if len(question.Options) > 0 {
		update["options"] = question.Options
	}
	if question.CorrectAnswer != "" {
		update["correct_answer"] = question.CorrectAnswer
	}
	if question.Status != "" {
		update["status"] = question.Status
	}

	if len(update) == 0 {
		return fmt.Errorf("no fields to update")
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": question.ID}, bson.M{"$set": update})
	return err
}

func (r *questionRepository) DeleteQuestion(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
