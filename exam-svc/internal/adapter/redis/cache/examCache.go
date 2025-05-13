package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/pkg/redis"
)

const examPrefix = "exam:%s"

type ExamCache struct {
	client *redis.Client
}

func NewExamCache(client *redis.Client) *ExamCache {
	return &ExamCache{client: client}
}

func (c *ExamCache) Set(ctx context.Context, exam domain.Exam) error {
	data, err := json.Marshal(exam)
	if err != nil {
		return err
	}
	return c.client.Unwrap().Set(ctx, fmt.Sprintf(examPrefix, exam.ID.Hex()), data, c.client.TTL()).Err()
}

func (c *ExamCache) Get(ctx context.Context, id string) (domain.Exam, error) {
	var exam domain.Exam
	data, err := c.client.Unwrap().Get(ctx, fmt.Sprintf(examPrefix, id)).Bytes()
	if err != nil {
		return exam, err
	}
	err = json.Unmarshal(data, &exam)
	return exam, err
}

func (c *ExamCache) Delete(ctx context.Context, id string) error {
	return c.client.Unwrap().Del(ctx, fmt.Sprintf(examPrefix, id)).Err()
}

func (c *ExamCache) SetMany(ctx context.Context, exams []domain.Exam) error {
	pipe := c.client.Unwrap().Pipeline()
	for _, e := range exams {
		data, err := json.Marshal(e)
		if err != nil {
			return err
		}
		pipe.Set(ctx, fmt.Sprintf(examPrefix, e.ID.Hex()), data, c.client.TTL())
	}
	_, err := pipe.Exec(ctx)
	return err
}
