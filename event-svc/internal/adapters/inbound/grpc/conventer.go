package grpc

import (
	"event-svc/internal/domain/model"

	eventsv1 "github.com/suyundykovv/margulan-protos/gen/go/events/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertCreateLessonRequest(req *eventsv1.CreateLessonRequest) (*model.Lesson, error) {
	return &model.Lesson{
		Title:      req.Title,
		StartTime:  req.StartTime.AsTime(),
		EndTime:    req.EndTime.AsTime(),
		GroupID:    req.GroupId,
		CourseID:   req.CourseId,
		MeetingURL: &req.MeetingUrl,
		Classroom:  req.Classroom,
		IsOnline:   req.IsOnline,
		Status:     model.LessonPlanned,
	}, nil
}

func convertLessonToProto(lesson *model.Lesson) *eventsv1.Lesson {
	return &eventsv1.Lesson{
		Id:         lesson.ID,
		Title:      lesson.Title,
		StartTime:  timestamppb.New(lesson.StartTime),
		EndTime:    timestamppb.New(lesson.EndTime),
		GroupId:    lesson.GroupID,
		CourseId:   lesson.CourseID,
		MeetingUrl: *lesson.MeetingURL,
		Classroom:  lesson.Classroom,
		IsOnline:   lesson.IsOnline,
		Status:     eventsv1.LessonStatus(lesson.Status),
	}
}
