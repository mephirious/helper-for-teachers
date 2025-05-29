package model

import (
	"fmt"
	"time"
)

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

func (t *Task) Validate() error {
	if t.Title == "" {
		return fmt.Errorf("title is required")
	}
	if t.DueDate.IsZero() {
		return fmt.Errorf("due date is required")
	}
	return nil
}
