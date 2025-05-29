package model

import (
	"fmt"
	"time"
)

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

func (s *LessonSchedule) Validate() error {
	if s.GroupID == "" {
		return fmt.Errorf("group ID is required")
	}
	if s.Title == "" {
		return fmt.Errorf("title is required")
	}
	if s.ValidFrom.IsZero() || s.ValidTo.IsZero() {
		return fmt.Errorf("valid dates are required")
	}
	if s.ValidFrom.After(s.ValidTo) {
		return fmt.Errorf("valid from must be before valid to")
	}
	return nil
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

func (s *TaskSchedule) Validate() error {
	if s.GroupID == "" {
		return fmt.Errorf("group ID is required")
	}
	if s.Title == "" {
		return fmt.Errorf("title is required")
	}
	if s.ValidFrom.IsZero() || s.ValidTo.IsZero() {
		return fmt.Errorf("valid dates are required")
	}
	if s.ValidFrom.After(s.ValidTo) {
		return fmt.Errorf("valid from must be before valid to")
	}
	return nil
}

type GroupSchedules struct {
	LessonSchedules []*LessonSchedule
	TaskSchedules   []*TaskSchedule
}
