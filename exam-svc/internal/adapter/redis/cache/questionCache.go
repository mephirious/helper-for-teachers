package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/pkg/redis"
)

const questionPrefix = "question:%s"

type QuestionCache struct {
	client *redis.Client
}

func NewQuestionCache(client *redis.Client) *QuestionCache {
	return &QuestionCache{client: client}
}

func (c *QuestionCache) Set(ctx context.Context, question domain.Question) error {
	data, err := json.Marshal(question)
	if err != nil {
		return err
	}
	return c.client.Unwrap().Set(ctx, fmt.Sprintf(questionPrefix, question.ID.Hex()), data, c.client.TTL()).Err()
}

func (c *QuestionCache) Get(ctx context.Context, id string) (domain.Question, error) {
	var question domain.Question
	data, err := c.client.Unwrap().Get(ctx, fmt.Sprintf(questionPrefix, id)).Bytes()
	if err != nil {
		return question, err
	}
	err = json.Unmarshal(data, &question)
	return question, err
}

func (c *QuestionCache) Delete(ctx context.Context, id string) error {
	return c.client.Unwrap().Del(ctx, fmt.Sprintf(questionPrefix, id)).Err()
}

func (c *QuestionCache) SetMany(ctx context.Context, questions []domain.Question) error {
	pipe := c.client.Unwrap().Pipeline()
	for _, q := range questions {
		data, err := json.Marshal(q)
		if err != nil {
			return err
		}
		pipe.Set(ctx, fmt.Sprintf(questionPrefix, q.ID.Hex()), data, c.client.TTL())
	}
	_, err := pipe.Exec(ctx)
	return err
}
