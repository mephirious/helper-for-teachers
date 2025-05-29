package model

import (
	"fmt"
	"time"
)

type Lesson struct {
	ID          string
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	GroupID     string
	CourseID    string
	Status      LessonStatus
	MeetingURL  *string
	Classroom   string
	IsOnline    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type LessonStatus int

const (
	LessonPlanned LessonStatus = iota
	LessonInProgress
	LessonCompleted
	LessonCanceled
)

func (s LessonStatus) String() string {
	return [...]string{"planned", "in_progress", "completed", "canceled"}[s]
}
func (l *Lesson) Validate() error {
	if l.Title == "" {
		return fmt.Errorf("title is required")
	}
	if l.StartTime.IsZero() || l.EndTime.IsZero() {
		return fmt.Errorf("start and end time are required")
	}
	if l.StartTime.After(l.EndTime) {
		return fmt.Errorf("start time must be before end time")
	}
	return nil
}
