package cache

import (
	"sync"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
)

type QuestionCache struct {
	data map[string]domain.Question
	mu   sync.RWMutex
}

func NewQuestionCache() *QuestionCache {
	return &QuestionCache{
		data: make(map[string]domain.Question),
	}
}

func (c *QuestionCache) Set(question domain.Question) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[question.ID.Hex()] = question
}

func (c *QuestionCache) SetMany(questions []domain.Question) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, q := range questions {
		c.data[q.ID.Hex()] = q
	}
}

func (c *QuestionCache) Get(id string) (domain.Question, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	question, ok := c.data[id]
	return question, ok
}

func (c *QuestionCache) GetAll() []domain.Question {
	c.mu.RLock()
	defer c.mu.RUnlock()
	questions := make([]domain.Question, 0, len(c.data))
	for _, q := range c.data {
		questions = append(questions, q)
	}
	return questions
}

func (c *QuestionCache) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, id)
}
