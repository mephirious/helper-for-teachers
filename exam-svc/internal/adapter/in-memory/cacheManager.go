package inmemory

import "github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/in-memory/cache"

type CacheManager struct {
	ExamCache     ExamCacheInterface
	QuestionCache QuestionCacheInterface
	TaskCache     TaskCacheInterface
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		ExamCache:     cache.NewExamCache(),
		QuestionCache: cache.NewQuestionCache(),
		TaskCache:     cache.NewTaskCache(),
	}
}
