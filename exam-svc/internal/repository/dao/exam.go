package dao

import (
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Exam struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	CreatedBy   primitive.ObjectID `bson:"created_by"`
	Status      string             `bson:"status"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

func FromDomainExam(e *domain.Exam) *Exam {
	return &Exam{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		CreatedBy:   e.CreatedBy,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func (e *Exam) ToDomainExam() *domain.Exam {
	return &domain.Exam{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		CreatedBy:   e.CreatedBy,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
