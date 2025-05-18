package model

import "time"

type Task struct {
	ID               string
	Title            string
	Description      string
	DueDate          time.Time
	GroupID          string
	CourseID         string
	Type             TaskType
	Status           TaskStatus
	ExternalResource string
	Attachments      []string
	MaxScore         *int32
	LessonID         *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type TaskType int

const (
	TaskExam TaskType = iota
	TaskAssignment
	TaskHomework
	TaskQuiz
	TaskProject
)

type TaskStatus int

const (
	TaskActive TaskStatus = iota
	TaskCompleted
	TaskGraded
	TaskArchived
)
