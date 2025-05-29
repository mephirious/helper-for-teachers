package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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
func (r *LessonRepository) GetByID(ctx context.Context, id string) (*model.Lesson, error) {
	if id == "" {
		return nil, fmt.Errorf("lesson ID cannot be empty")
	}
	query := `
        SELECT 
            id, 
            title, 
            start_time, 
            end_time, 
            group_id, 
            course_id, 
            status, 
            meeting_url, 
            classroom, 
            is_online, 
            created_at, 
            updated_at
        FROM lessons 
        WHERE id = $1
    `
	row := r.db.QueryRowContext(ctx, query, id)
	var lesson model.Lesson
	err := row.Scan(
		&lesson.ID,
		&lesson.Title,
		&lesson.StartTime,
		&lesson.EndTime,
		&lesson.GroupID,
		&lesson.CourseID,
		&lesson.Status,
		&lesson.MeetingURL,
		&lesson.Classroom,
		&lesson.IsOnline,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("lesson not found")
		}
		return nil, fmt.Errorf("failed to get lesson: %w", err)
	}

	return &lesson, nil
}

func (r *LessonRepository) GetAll(ctx context.Context) (*[]model.Lesson, error) {
	query := `
        SELECT 
            id, 
            title, 
            start_time, 
            end_time, 
            group_id, 
            course_id, 
            status, 
            meeting_url, 
            classroom, 
            is_online, 
            created_at, 
            updated_at
        FROM lessons 
        WHERE id = $1
    `
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var lessons []*model.Lesson
	for rows.Next() {
		var lesson model.Lesson
		err := row.Scan(
			&lesson.ID,
			&lesson.Title,
			&lesson.StartTime,
			&lesson.EndTime,
			&lesson.GroupID,
			&lesson.CourseID,
			&lesson.Status,
			&lesson.MeetingURL,
			&lesson.Classroom,
			&lesson.IsOnline,
			&lesson.CreatedAt,
			&lesson.UpdatedAt,
		); err != nil{
		return nil, err}
		lessons = append(lessons,&lesson)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	} 
	return lessons, err
}
