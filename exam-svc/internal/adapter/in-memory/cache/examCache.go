package cache

import (
	"sync"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
)

type ExamCache struct {
	data map[string]domain.Exam
	mu   sync.RWMutex
}

func NewExamCache() *ExamCache {
	return &ExamCache{
		data: make(map[string]domain.Exam),
	}
}

func (c *ExamCache) Set(exam domain.Exam) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[exam.ID.Hex()] = exam
}

func (c *ExamCache) SetMany(exams []domain.Exam) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, exam := range exams {
		c.data[exam.ID.Hex()] = exam
	}
}

func (c *ExamCache) Get(id string) (domain.Exam, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	exam, ok := c.data[id]
	return exam, ok
}

func (c *ExamCache) GetAll() []domain.Exam {
	c.mu.RLock()
	defer c.mu.RUnlock()
	exams := make([]domain.Exam, 0, len(c.data))
	for _, exam := range c.data {
		exams = append(exams, exam)
	}
	return exams
}

func (c *ExamCache) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, id)
}
