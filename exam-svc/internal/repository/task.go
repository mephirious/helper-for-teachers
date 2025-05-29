package repository

import (
	"context"
	"fmt"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/repository/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

func NewTaskRepository(db *mongo.Database, client *mongo.Client) TaskRepository {
	return &taskRepository{
		collection: db.Collection("tasks"),
		client:     client,
	}
}

func (r *taskRepository) CreateTask(ctx context.Context, task *domain.Task) error {
	session, err := r.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
		examColl := r.collection.Database().Collection("exams")
		var examDAO dao.Exam
		err := examColl.FindOne(sessionCtx, bson.M{"_id": task.ExamID}).Decode(&examDAO)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("exam with id %s not found", task.ExamID.Hex())
			}
			return nil, fmt.Errorf("failed to verify exam: %w", err)
		}

		taskDAO := dao.FromDomainTask(task)
		_, err = r.collection.InsertOne(sessionCtx, taskDAO)
		if err != nil {
			return nil, fmt.Errorf("failed to insert task: %w", err)
		}

		return nil, nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id primitive.ObjectID) (*domain.Task, error) {
	var taskDAO dao.Task
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&taskDAO)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return taskDAO.ToDomainTask(), nil
}

func (r *taskRepository) GetTasksByExamID(ctx context.Context, examID primitive.ObjectID) ([]domain.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"exam_id": examID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var taskDAOs []dao.Task
	if err := cursor.All(ctx, &taskDAOs); err != nil {
		return nil, err
	}

	tasks := make([]domain.Task, 0, len(taskDAOs))
	for _, daoTask := range taskDAOs {
		tasks = append(tasks, *daoTask.ToDomainTask())
	}
	return tasks, nil
}

func (r *taskRepository) GetAllTasks(ctx context.Context) ([]domain.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var taskDAOs []dao.Task
	if err := cursor.All(ctx, &taskDAOs); err != nil {
		return nil, err
	}

	tasks := make([]domain.Task, 0, len(taskDAOs))
	for _, daoTask := range taskDAOs {
		tasks = append(tasks, *daoTask.ToDomainTask())
	}
	return tasks, nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *domain.Task) error {
	update := bson.M{}
	if task.TaskType != "" {
		update["task_type"] = task.TaskType
	}
	if task.Description != "" {
		update["description"] = task.Description
	}
	if task.Score != 0 {
		update["score"] = task.Score
	}

	if len(update) == 0 {
		return fmt.Errorf("no fields to update")
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": task.ID}, bson.M{"$set": update})
	return err
}

func (r *taskRepository) DeleteTask(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
