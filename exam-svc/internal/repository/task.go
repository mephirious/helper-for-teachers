package repository

import (
	"context"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/repository/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) TaskRepository {
	return &taskRepository{
		collection: db.Collection("tasks"),
	}
}

func (r *taskRepository) CreateTask(ctx context.Context, task *domain.Task) error {
	taskDAO := dao.FromDomainTask(task)
	_, err := r.collection.InsertOne(ctx, taskDAO)
	return err
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

func (r *taskRepository) DeleteTask(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
