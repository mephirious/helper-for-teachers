package cache

import (
	"sync"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
)

type TaskCache struct {
	data map[string]domain.Task
	mu   sync.RWMutex
}

func NewTaskCache() *TaskCache {
	return &TaskCache{
		data: make(map[string]domain.Task),
	}
}

func (c *TaskCache) Set(task domain.Task) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[task.ID.Hex()] = task
}

func (c *TaskCache) SetMany(tasks []domain.Task) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, task := range tasks {
		c.data[task.ID.Hex()] = task
	}
}

func (c *TaskCache) Get(id string) (domain.Task, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	task, ok := c.data[id]
	return task, ok
}

func (c *TaskCache) GetAll() []domain.Task {
	c.mu.RLock()
	defer c.mu.RUnlock()
	tasks := make([]domain.Task, 0, len(c.data))
	for _, t := range c.data {
		tasks = append(tasks, t)
	}
	return tasks
}

func (c *TaskCache) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, id)
}
