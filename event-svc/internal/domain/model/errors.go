package model

import "errors"

var (
	// Common errors
	ErrInvalidID     = errors.New("invalid ID")
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")

	// Lesson-specific errors
	ErrInvalidLessonTime = errors.New("invalid lesson time range")

	// Task-specific errors
	ErrInvalidDueDate = errors.New("invalid due date")

	// Schedule-specific errors
	ErrInvalidScheduleRange = errors.New("invalid schedule date range")
)
