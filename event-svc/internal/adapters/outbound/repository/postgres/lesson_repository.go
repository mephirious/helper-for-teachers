package postgres

import (
	"context"
	"database/sql"

	"event-svc/internal/domain/model"
)

type LessonRepository struct {
	db *sql.DB
}

func NewLessonRepository(db *sql.DB) *LessonRepository {
	return &LessonRepository{db: db}
}

func (r *LessonRepository) Create(ctx context.Context, lesson *model.Lesson) error {
	query := `INSERT INTO lessons (id, title, start_time, end_time, group_id, course_id, 
		status, meeting_url, classroom, is_online, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.ExecContext(ctx, query,
		lesson.ID,
		lesson.Title,
		lesson.StartTime,
		lesson.EndTime,
		lesson.GroupID,
		lesson.CourseID,
		lesson.Status,
		lesson.MeetingURL,
		lesson.Classroom,
		lesson.IsOnline,
		lesson.CreatedAt,
		lesson.UpdatedAt,
	)

	return err
}
