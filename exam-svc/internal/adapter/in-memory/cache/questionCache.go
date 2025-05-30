package cache

import (
	"context"
	"sync"
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/repository"
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

func (c *QuestionCache) Init(ctx context.Context, repo repository.QuestionRepository) error {
	if err := c.refresh(ctx, repo); err != nil {
		return err
	}

	go c.startRefreshLoop(ctx, repo)
	return nil
}

func (c *QuestionCache) refresh(ctx context.Context, repo repository.QuestionRepository) error {
	questions, err := repo.GetAllQuestions(ctx)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]domain.Question, len(questions))
	for _, question := range questions {
		c.data[question.ID.Hex()] = question
	}
	return nil
}

func (c *QuestionCache) startRefreshLoop(ctx context.Context, repo repository.QuestionRepository) {
	ticker := time.NewTicker(12 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := c.refresh(ctx, repo); err != nil {
				println("Failed to refresh question cache:", err.Error())
			}
		case <-ctx.Done():
			return
		}
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
