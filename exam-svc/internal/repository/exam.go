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

type examRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

func NewExamRepository(db *mongo.Database, client *mongo.Client) ExamRepository {
	return &examRepository{
		collection: db.Collection("exams"),
		client:     client,
	}
}

func (r *examRepository) CreateExam(ctx context.Context, exam *domain.Exam) error {
	exam.ID = primitive.NewObjectID()
	exam.CreatedAt = time.Now()
	exam.UpdatedAt = exam.CreatedAt

	_, err := r.collection.InsertOne(ctx, dao.FromDomainExam(exam))
	return err
}

func (r *examRepository) GetExamByID(ctx context.Context, id primitive.ObjectID) (*domain.Exam, error) {
	var examDAO dao.Exam
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&examDAO)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return examDAO.ToDomainExam(), nil
}

func (r *examRepository) GetExamsByUser(ctx context.Context, userID primitive.ObjectID) ([]domain.Exam, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"created_by": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var examDAOs []dao.Exam
	if err := cursor.All(ctx, &examDAOs); err != nil {
		return nil, err
	}

	exams := make([]domain.Exam, 0, len(examDAOs))
	for _, e := range examDAOs {
		exams = append(exams, *e.ToDomainExam())
	}
	return exams, nil
}

func (r *examRepository) UpdateExamStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{"status": status, "updated_at": time.Now()},
	})
	return err
}

func (r *examRepository) UpdateExam(ctx context.Context, exam *domain.Exam) error {
	update := bson.M{}
	if exam.Title != "" {
		update["title"] = exam.Title
	}
	if exam.Description != "" {
		update["description"] = exam.Description
	}
	if exam.Status != "" {
		update["status"] = exam.Status
	}
	update["updated_at"] = time.Now()

	if len(update) == 0 {
		return fmt.Errorf("no fields to update")
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": exam.ID}, bson.M{"$set": update})
	return err
}

func (r *examRepository) DeleteExam(ctx context.Context, id primitive.ObjectID) error {
	session, err := r.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
		_, err := r.collection.DeleteOne(sessionCtx, bson.M{"_id": id})
		if err != nil {
			return nil, fmt.Errorf("failed to delete exam: %w", err)
		}

		taskColl := r.collection.Database().Collection("tasks")
		_, err = taskColl.DeleteMany(sessionCtx, bson.M{"exam_id": id})
		if err != nil {
			return nil, fmt.Errorf("failed to delete tasks: %w", err)
		}

		questionColl := r.collection.Database().Collection("questions")
		_, err = questionColl.DeleteMany(sessionCtx, bson.M{"exam_id": id})
		if err != nil {
			return nil, fmt.Errorf("failed to delete questions: %w", err)
		}

		return nil, nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}

func (r *examRepository) GetAllExams(ctx context.Context) ([]domain.Exam, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var examDAOs []dao.Exam
	if err := cursor.All(ctx, &examDAOs); err != nil {
		return nil, err
	}

	exams := make([]domain.Exam, 0, len(examDAOs))
	for _, e := range examDAOs {
		exams = append(exams, *e.ToDomainExam())
	}
	return exams, nil
}
