package model

import "time"

type LessonSchedule struct {
	ID        string
	GroupID   string
	Title     string
	ValidFrom time.Time
	ValidTo   time.Time
	IsActive  bool
	CourseID  string
	LessonIDs []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TaskSchedule struct {
	ID        string
	GroupID   string
	Title     string
	ValidFrom time.Time
	ValidTo   time.Time
	IsActive  bool
	CourseID  string
	TaskIDs   []string
	CreatedAt time.Time
	UpdatedAt time.Time
}
