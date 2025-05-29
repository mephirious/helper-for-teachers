package grpc

import (
	"event-svc/internal/domain/model"

	eventsv1 "github.com/suyundykovv/margulan-protos/gen/go/events/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Helper functions for pointer handling
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}

// Lesson converters
func ConvertCreateLessonRequestToDomain(req *eventsv1.CreateLessonRequest) *model.Lesson {
	var meetingURL *string
	if req.GetMeetingUrl() != nil {
		meetingURL = stringPtr(req.GetMeetingUrl().GetValue())
	}

	return &model.Lesson{
		Title:      req.GetTitle(),
		StartTime:  req.GetStartTime().AsTime(),
		EndTime:    req.GetEndTime().AsTime(),
		GroupID:    req.GetGroupId(),
		CourseID:   req.GetCourseId(),
		Status:     model.LessonStatus(req.GetStatus()),
		MeetingURL: meetingURL,
		Classroom:  req.GetClassroom(),
		IsOnline:   req.GetIsOnline(),
	}
}

func ConvertDomainLessonToProto(lesson *model.Lesson) *eventsv1.Lesson {
	var meetingURL *wrapperspb.StringValue
	if lesson.MeetingURL != nil {
		meetingURL = wrapperspb.String(*lesson.MeetingURL)
	}

	return &eventsv1.Lesson{
		Id:         lesson.ID,
		Title:      lesson.Title,
		StartTime:  timestamppb.New(lesson.StartTime),
		EndTime:    timestamppb.New(lesson.EndTime),
		GroupId:    lesson.GroupID,
		CourseId:   lesson.CourseID,
		Status:     eventsv1.LessonStatus(lesson.Status),
		MeetingUrl: meetingURL,
		Classroom:  lesson.Classroom,
		IsOnline:   lesson.IsOnline,
		CreatedAt:  timestamppb.New(lesson.CreatedAt),
		UpdatedAt:  timestamppb.New(lesson.UpdatedAt),
	}
}

// Schedule converters
func ConvertCreateScheduleRequestToDomain(req *eventsv1.CreateLessonScheduleRequest) *model.LessonSchedule {
	return &model.LessonSchedule{
		GroupID:   req.GetGroupId(),
		Title:     req.GetTitle(),
		ValidFrom: req.GetValidFrom().AsTime(),
		ValidTo:   req.GetValidTo().AsTime(),
		CourseID:  req.GetCourseId(),
		IsActive:  req.GetIsActive(),
		LessonIDs: req.GetLessonIds(),
	}
}

func ConvertDomainScheduleToProto(schedule *model.LessonSchedule) *eventsv1.LessonSchedule {
	return &eventsv1.LessonSchedule{
		Id:        schedule.ID,
		GroupId:   schedule.GroupID,
		Title:     schedule.Title,
		ValidFrom: timestamppb.New(schedule.ValidFrom),
		ValidTo:   timestamppb.New(schedule.ValidTo),
		CourseId:  schedule.CourseID,
		IsActive:  schedule.IsActive,
		LessonIds: schedule.LessonIDs,
		CreatedAt: timestamppb.New(schedule.CreatedAt),
		UpdatedAt: timestamppb.New(schedule.UpdatedAt),
	}
}

// Task converters
func ConvertCreateTaskRequestToDomain(req *eventsv1.CreateTaskRequest) *model.Task {

	var maxScore *int32
	if req.GetMaxScore() != nil {
		maxScore = int32Ptr(req.GetMaxScore().GetValue())
	}

	return &model.Task{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		DueDate:     req.GetDueDate().AsTime(),
		GroupID:     req.GetGroupId(),
		CourseID:    req.GetCourseId(),
		Type:        model.TaskType(req.GetType()),
		Attachments: req.GetAttachments(),
		MaxScore:    maxScore,
	}
}

func ConvertDomainTaskToProto(task *model.Task) *eventsv1.Task {

	var maxScore *wrapperspb.Int32Value
	if task.MaxScore != nil {
		maxScore = wrapperspb.Int32(*task.MaxScore)
	}

	var lessonID *wrapperspb.StringValue
	if task.LessonID != nil {
		lessonID = wrapperspb.String(*task.LessonID)
	}

	return &eventsv1.Task{
		Id:                 task.ID,
		Title:              task.Title,
		Description:        task.Description,
		DueDate:            timestamppb.New(task.DueDate),
		GroupId:            task.GroupID,
		CourseId:           task.CourseID,
		Type:               eventsv1.TaskType(task.Type),
		Status:             eventsv1.TaskStatus(task.Status),
		ExternalResourceId: task.ExternalResource,
		Attachments:        task.Attachments,
		MaxScore:           maxScore,
		LessonId:           lessonID,
		CreatedAt:          timestamppb.New(task.CreatedAt),
		UpdatedAt:          timestamppb.New(task.UpdatedAt),
	}
}
