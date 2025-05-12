package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/pkg/redis"
)

const taskPrefix = "task:%s"

type TaskCache struct {
	client *redis.Client
}

func NewTaskCache(client *redis.Client) *TaskCache {
	return &TaskCache{client: client}
}

func (c *TaskCache) Set(ctx context.Context, task domain.Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return c.client.Unwrap().Set(ctx, fmt.Sprintf(taskPrefix, task.ID.Hex()), data, 0).Err()
}

func (c *TaskCache) Get(ctx context.Context, id string) (domain.Task, error) {
	var task domain.Task
	data, err := c.client.Unwrap().Get(ctx, fmt.Sprintf(taskPrefix, id)).Bytes()
	if err != nil {
		return task, err
	}
	err = json.Unmarshal(data, &task)
	return task, err
}

func (c *TaskCache) Delete(ctx context.Context, id string) error {
	return c.client.Unwrap().Del(ctx, fmt.Sprintf(taskPrefix, id)).Err()
}

func (c *TaskCache) SetMany(ctx context.Context, tasks []domain.Task) error {
	pipe := c.client.Unwrap().Pipeline()
	for _, t := range tasks {
		data, err := json.Marshal(t)
		if err != nil {
			return err
		}
		pipe.Set(ctx, fmt.Sprintf(taskPrefix, t.ID.Hex()), data, 0)
	}
	_, err := pipe.Exec(ctx)
	return err
}
