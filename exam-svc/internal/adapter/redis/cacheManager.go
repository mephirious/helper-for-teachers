package redis

import (
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/redis/cache"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/pkg/redis"
)

type CacheManager struct {
	ExamCache     ExamCache
	QuestionCache QuestionCache
	TaskCache     TaskCache
}

func NewCacheManager(redisClient *redis.Client) *CacheManager {
	return &CacheManager{
		ExamCache:     cache.NewExamCache(redisClient),
		QuestionCache: cache.NewQuestionCache(redisClient),
		TaskCache:     cache.NewTaskCache(redisClient),
	}
}
