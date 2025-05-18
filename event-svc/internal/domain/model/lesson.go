package model

import "time"

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
