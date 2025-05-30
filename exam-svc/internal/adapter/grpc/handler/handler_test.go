package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/usecase/mocks"
	pb "github.com/mephirious/helper-for-teachers/services/exam-svc/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestExamHandler_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("invalid exam ID", func(t *testing.T) {
		req := &pb.CreateTaskRequest{
			Task: &pb.Task{
				ExamId: "invalid",
			},
		}

		resp, err := handler.CreateTask(ctx, req)
		assert.Error(t, err, "expected error for invalid exam ID")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "error code should be InvalidArgument")
		assert.Contains(t, st.Message(), "invalid exam_id", "error message should mention invalid exam_id")
	})
}

func TestExamHandler_GetTaskByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("invalid task ID", func(t *testing.T) {
		resp, err := handler.GetTaskByID(ctx, &pb.GetTaskByIDRequest{Id: "invalid"})
		assert.Error(t, err, "expected error for invalid task ID")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "error code should be InvalidArgument")
		assert.Contains(t, st.Message(), "invalid id", "error message should mention invalid id")
	})

	t.Run("task not found", func(t *testing.T) {
		taskID := primitive.NewObjectID()
		taskUC.EXPECT().GetTaskByID(ctx, taskID).Return(nil, errors.New("not found"))

		resp, err := handler.GetTaskByID(ctx, &pb.GetTaskByIDRequest{Id: taskID.Hex()})
		assert.Error(t, err, "expected error when task not found")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.NotFound, st.Code(), "error code should be NotFound")
		assert.Contains(t, st.Message(), "task not found", "error message should mention not found")
	})
}

func TestExamHandler_GetTasksByExamID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("successful tasks retrieval", func(t *testing.T) {
		examID := primitive.NewObjectID()
		createdAt := time.Now()
		tasks := []domain.Task{
			{
				ID:          primitive.NewObjectID(),
				ExamID:      examID,
				TaskType:    "essay",
				Description: "Write an essay",
				Score:       10,
				CreatedAt:   createdAt,
			},
			{
				ID:          primitive.NewObjectID(),
				ExamID:      examID,
				TaskType:    "multiple_choice",
				Description: "Choose the correct answer",
				Score:       5,
				CreatedAt:   createdAt,
			},
		}

		taskUC.EXPECT().GetTasksByExamID(ctx, examID).Return(tasks, nil)

		resp, err := handler.GetTasksByExamID(ctx, &pb.GetTasksByExamIDRequest{ExamId: examID.Hex()})
		require.NoError(t, err, "expected no error when retrieving tasks")
		assert.Len(t, resp.Tasks, 2, "should return two tasks")
		assert.Equal(t, tasks[0].ID.Hex(), resp.Tasks[0].Id, "first task ID should match")
		assert.Equal(t, tasks[1].Description, resp.Tasks[1].Description, "second task description should match")
	})

	t.Run("invalid exam ID", func(t *testing.T) {
		resp, err := handler.GetTasksByExamID(ctx, &pb.GetTasksByExamIDRequest{ExamId: "invalid"})
		assert.Error(t, err, "expected error for invalid exam ID")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "error code should be InvalidArgument")
	})

	t.Run("use case error", func(t *testing.T) {
		examID := primitive.NewObjectID()
		taskUC.EXPECT().GetTasksByExamID(ctx, examID).Return(nil, errors.New("database error"))

		resp, err := handler.GetTasksByExamID(ctx, &pb.GetTasksByExamIDRequest{ExamId: examID.Hex()})
		assert.Error(t, err, "expected error from use case")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.Internal, st.Code(), "error code should be Internal")
	})
}

func TestExamHandler_GetAllTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("successful all tasks retrieval", func(t *testing.T) {
		examID := primitive.NewObjectID()
		createdAt := time.Now()
		tasks := []domain.Task{
			{
				ID:          primitive.NewObjectID(),
				ExamID:      examID,
				TaskType:    "essay",
				Description: "Write an essay",
				Score:       10,
				CreatedAt:   createdAt,
			},
		}

		taskUC.EXPECT().GetAllTasks(ctx).Return(tasks, nil)

		resp, err := handler.GetAllTasks(ctx, &emptypb.Empty{})
		require.NoError(t, err, "expected no error when retrieving all tasks")
		assert.Len(t, resp.Tasks, 1, "should return one task")
		assert.Equal(t, tasks[0].ID.Hex(), resp.Tasks[0].Id, "task ID should match")
	})

	t.Run("use case error", func(t *testing.T) {
		taskUC.EXPECT().GetAllTasks(ctx).Return(nil, errors.New("database error"))

		resp, err := handler.GetAllTasks(ctx, &emptypb.Empty{})
		assert.Error(t, err, "expected error from use case")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.Internal, st.Code(), "error code should be Internal")
	})
}

func TestExamHandler_UpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("invalid task ID", func(t *testing.T) {
		req := &pb.UpdateTaskRequest{
			Task: &pb.Task{Id: "invalid"},
		}

		resp, err := handler.UpdateTask(ctx, req)
		assert.Error(t, err, "expected error for invalid task ID")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "error code should be InvalidArgument")
	})
}

func TestExamHandler_DeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("successful task deletion", func(t *testing.T) {
		taskID := primitive.NewObjectID()
		taskUC.EXPECT().DeleteTask(ctx, taskID).Return(nil)

		resp, err := handler.DeleteTask(ctx, &pb.DeleteTaskRequest{Id: taskID.Hex()})
		require.NoError(t, err, "expected no error when deleting task")
		assert.NotNil(t, resp, "response should be non-nil")
	})

	t.Run("invalid task ID", func(t *testing.T) {
		resp, err := handler.DeleteTask(ctx, &pb.DeleteTaskRequest{Id: "invalid"})
		assert.Error(t, err, "expected error for invalid task ID")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "error code should be InvalidArgument")
	})
}

func TestExamHandler_CreateQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("invalid exam ID", func(t *testing.T) {
		req := &pb.CreateQuestionRequest{
			Question: &pb.Question{ExamId: "invalid"},
		}

		resp, err := handler.CreateQuestion(ctx, req)
		assert.Error(t, err, "expected error for invalid exam ID")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "error code should be InvalidArgument")
	})
}

func TestExamHandler_GetQuestionByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("successful question retrieval", func(t *testing.T) {
		questionID := primitive.NewObjectID()
		examID := primitive.NewObjectID()
		createdAt := time.Now()
		question := &domain.Question{
			ID:            questionID,
			ExamID:        examID,
			QuestionText:  "What is 2+2?",
			Options:       []string{"2", "4", "6"},
			CorrectAnswer: "4",
			Status:        "active",
			CreatedAt:     createdAt,
		}

		questionUC.EXPECT().GetQuestionByID(ctx, questionID).Return(question, nil)

		resp, err := handler.GetQuestionByID(ctx, &pb.GetQuestionByIDRequest{Id: questionID.Hex()})
		require.NoError(t, err, "expected no error when retrieving question")
		assert.Equal(t, questionID.Hex(), resp.Id, "question ID should match")
		assert.Equal(t, examID.Hex(), resp.ExamId, "exam ID should match")
		assert.Equal(t, "What is 2+2?", resp.QuestionText, "question text should match")
		assert.Equal(t, []string{"2", "4", "6"}, resp.Options, "options should match")
		assert.Equal(t, "4", resp.CorrectAnswer, "correct answer should match")
		assert.Equal(t, "active", resp.Status, "status should match")
		assert.WithinDuration(t, createdAt, resp.CreatedAt.AsTime(), time.Second, "created at should match")
	})

	t.Run("invalid question ID", func(t *testing.T) {
		resp, err := handler.GetQuestionByID(ctx, &pb.GetQuestionByIDRequest{Id: "invalid"})
		assert.Error(t, err, "expected error for invalid question ID")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "error code should be InvalidArgument")
	})

	t.Run("question not found", func(t *testing.T) {
		questionID := primitive.NewObjectID()
		questionUC.EXPECT().GetQuestionByID(ctx, questionID).Return(nil, errors.New("not found"))

		resp, err := handler.GetQuestionByID(ctx, &pb.GetQuestionByIDRequest{Id: questionID.Hex()})
		assert.Error(t, err, "expected error when question not found")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.NotFound, st.Code(), "error code should be NotFound")
	})
}

func TestExamHandler_GetQuestionsByExamID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("successful questions retrieval", func(t *testing.T) {
		examID := primitive.NewObjectID()
		createdAt := time.Now()
		questions := []domain.Question{
			{
				ID:            primitive.NewObjectID(),
				ExamID:        examID,
				QuestionText:  "What is 2+2?",
				Options:       []string{"2", "4", "6"},
				CorrectAnswer: "4",
				Status:        "active",
				CreatedAt:     createdAt,
			},
		}

		questionUC.EXPECT().GetQuestionsByExamID(ctx, examID).Return(questions, nil)

		resp, err := handler.GetQuestionsByExamID(ctx, &pb.GetQuestionsByExamIDRequest{ExamId: examID.Hex()})
		require.NoError(t, err, "expected no error when retrieving questions")
		assert.Len(t, resp.Questions, 1, "should return one question")
		assert.Equal(t, questions[0].ID.Hex(), resp.Questions[0].Id, "question ID should match")
	})

	t.Run("invalid exam ID", func(t *testing.T) {
		resp, err := handler.GetQuestionsByExamID(ctx, &pb.GetQuestionsByExamIDRequest{ExamId: "invalid"})
		assert.Error(t, err, "expected error for invalid exam ID")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "error code should be InvalidArgument")
	})
}

func TestExamHandler_GetAllQuestions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("successful all questions retrieval", func(t *testing.T) {
		examID := primitive.NewObjectID()
		createdAt := time.Now()
		questions := []domain.Question{
			{
				ID:            primitive.NewObjectID(),
				ExamID:        examID,
				QuestionText:  "What is 2+2?",
				Options:       []string{"2", "4", "6"},
				CorrectAnswer: "4",
				Status:        "active",
				CreatedAt:     createdAt,
			},
		}

		questionUC.EXPECT().GetAllQuestions(ctx).Return(questions, nil)

		resp, err := handler.GetAllQuestions(ctx, &emptypb.Empty{})
		require.NoError(t, err, "expected no error when retrieving all questions")
		assert.Len(t, resp.Questions, 1, "should return one question")
		assert.Equal(t, questions[0].ID.Hex(), resp.Questions[0].Id, "question ID should match")
	})

	t.Run("use case error", func(t *testing.T) {
		questionUC.EXPECT().GetAllQuestions(ctx).Return(nil, errors.New("database error"))

		resp, err := handler.GetAllQuestions(ctx, &emptypb.Empty{})
		assert.Error(t, err, "expected error from use case")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.Internal, st.Code(), "error code should be Internal")
	})
}

func TestExamHandler_UpdateQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("invalid question ID", func(t *testing.T) {
		req := &pb.UpdateQuestionRequest{
			Question: &pb.Question{Id: "invalid"},
		}

		resp, err := handler.UpdateQuestion(ctx, req)
		assert.Error(t, err, "expected error for invalid question ID")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "error code should be InvalidArgument")
	})
}

func TestExamHandler_DeleteQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("successful question deletion", func(t *testing.T) {
		questionID := primitive.NewObjectID()
		questionUC.EXPECT().DeleteQuestion(ctx, questionID).Return(nil)

		resp, err := handler.DeleteQuestion(ctx, &pb.DeleteQuestionRequest{Id: questionID.Hex()})
		require.NoError(t, err, "expected no error when deleting question")
		assert.NotNil(t, resp, "response should be non-nil")
	})

	t.Run("invalid question ID", func(t *testing.T) {
		resp, err := handler.DeleteQuestion(ctx, &pb.DeleteQuestionRequest{Id: "invalid"})
		assert.Error(t, err, "expected error for invalid question ID")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "error code should be InvalidArgument")
	})
}

func TestExamHandler_CreateExam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("invalid created by ID", func(t *testing.T) {
		req := &pb.CreateExamRequest{
			Exam: &pb.Exam{CreatedBy: "invalid"},
		}

		resp, err := handler.CreateExam(ctx, req)
		assert.Error(t, err, "expected error for invalid created by ID")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "error code should be InvalidArgument")
	})
}

func TestExamHandler_GetExamByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("successful exam retrieval", func(t *testing.T) {
		examID := primitive.NewObjectID()
		userID := primitive.NewObjectID()
		createdAt := time.Now()
		exam := &domain.Exam{
			ID:          examID,
			Title:       "Math Exam",
			Description: "Basic math exam",
			CreatedBy:   userID,
			Status:      "draft",
			CreatedAt:   createdAt,
			UpdatedAt:   createdAt,
		}

		examUC.EXPECT().GetExamByID(ctx, examID).Return(exam, nil)

		resp, err := handler.GetExamByID(ctx, &pb.GetExamByIDRequest{Id: examID.Hex()})
		require.NoError(t, err, "expected no error when retrieving exam")
		assert.Equal(t, examID.Hex(), resp.Id, "exam ID should match")
		assert.Equal(t, "Math Exam", resp.Title, "title should match")
		assert.Equal(t, userID.Hex(), resp.CreatedBy, "created by should match")
	})

	t.Run("invalid exam ID", func(t *testing.T) {
		resp, err := handler.GetExamByID(ctx, &pb.GetExamByIDRequest{Id: "invalid"})
		assert.Error(t, err, "expected error for invalid exam ID")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "error code should be InvalidArgument")
	})

	t.Run("exam not found", func(t *testing.T) {
		examID := primitive.NewObjectID()
		examUC.EXPECT().GetExamByID(ctx, examID).Return(nil, errors.New("not found"))

		resp, err := handler.GetExamByID(ctx, &pb.GetExamByIDRequest{Id: examID.Hex()})
		assert.Error(t, err, "expected error when exam not found")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.NotFound, st.Code(), "error code should be NotFound")
	})
}

func TestExamHandler_GetExamsByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("successful exams retrieval", func(t *testing.T) {
		userID := primitive.NewObjectID()
		createdAt := time.Now()
		exams := []domain.Exam{
			{
				ID:          primitive.NewObjectID(),
				Title:       "Math Exam",
				Description: "Basic math exam",
				CreatedBy:   userID,
				Status:      "draft",
				CreatedAt:   createdAt,
				UpdatedAt:   createdAt,
			},
		}

		examUC.EXPECT().GetExamsByUser(ctx, userID).Return(exams, nil)

		resp, err := handler.GetExamsByUser(ctx, &pb.GetExamsByUserRequest{UserId: userID.Hex()})
		require.NoError(t, err, "expected no error when retrieving exams")
		assert.Len(t, resp.Exams, 1, "should return one exam")
		assert.Equal(t, exams[0].ID.Hex(), resp.Exams[0].Id, "exam ID should match")
	})

	t.Run("invalid user ID", func(t *testing.T) {
		resp, err := handler.GetExamsByUser(ctx, &pb.GetExamsByUserRequest{UserId: "invalid"})
		assert.Error(t, err, "expected error for invalid user ID")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "error code should be InvalidArgument")
	})
}

func TestExamHandler_UpdateExam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("invalid exam ID", func(t *testing.T) {
		req := &pb.UpdateExamRequest{
			Exam: &pb.Exam{Id: "invalid"},
		}

		resp, err := handler.UpdateExam(ctx, req)
		assert.Error(t, err, "expected error for invalid exam ID")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "error code should be InvalidArgument")
	})
}

func TestExamHandler_UpdateExamStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("successful exam status update", func(t *testing.T) {
		examID := primitive.NewObjectID()
		examUC.EXPECT().UpdateExamStatus(ctx, examID, "published").Return(nil)

		resp, err := handler.UpdateExamStatus(ctx, &pb.UpdateExamStatusRequest{Id: examID.Hex(), Status: "published"})
		require.NoError(t, err, "expected no error when updating exam status")
		assert.NotNil(t, resp, "response should be non-nil")
	})

	t.Run("invalid exam ID", func(t *testing.T) {
		resp, err := handler.UpdateExamStatus(ctx, &pb.UpdateExamStatusRequest{Id: "invalid"})
		assert.Error(t, err, "expected error for invalid exam ID")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "error code should be InvalidArgument")
	})
}

func TestExamHandler_DeleteExam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("successful exam deletion", func(t *testing.T) {
		examID := primitive.NewObjectID()
		examUC.EXPECT().DeleteExam(ctx, examID).Return(nil)

		resp, err := handler.DeleteExam(ctx, &pb.DeleteExamRequest{Id: examID.Hex()})
		require.NoError(t, err, "expected no error when deleting exam")
		assert.NotNil(t, resp, "response should be non-nil")
	})

	t.Run("invalid exam ID", func(t *testing.T) {
		resp, err := handler.DeleteExam(ctx, &pb.DeleteExamRequest{Id: "invalid"})
		assert.Error(t, err, "expected error for invalid exam ID")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "error code should be InvalidArgument")
	})
}

func TestExamHandler_GetAllExams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("successful all exams retrieval", func(t *testing.T) {
		userID := primitive.NewObjectID()
		createdAt := time.Now()
		exams := []domain.Exam{
			{
				ID:          primitive.NewObjectID(),
				Title:       "Math Exam",
				Description: "Basic math exam",
				CreatedBy:   userID,
				Status:      "draft",
				CreatedAt:   createdAt,
				UpdatedAt:   createdAt,
			},
		}

		examUC.EXPECT().GetAllExams(ctx).Return(exams, nil)

		resp, err := handler.GetAllExams(ctx, &emptypb.Empty{})
		require.NoError(t, err, "expected no error when retrieving all exams")
		assert.Len(t, resp.Exams, 1, "should return one exam")
		assert.Equal(t, exams[0].ID.Hex(), resp.Exams[0].Id, "exam ID should match")
	})

	t.Run("use case error", func(t *testing.T) {
		examUC.EXPECT().GetAllExams(ctx).Return(nil, errors.New("database error"))

		resp, err := handler.GetAllExams(ctx, &emptypb.Empty{})
		assert.Error(t, err, "expected error from use case")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.Internal, st.Code(), "error code should be Internal")
	})
}

func TestExamHandler_GetExamWithDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("invalid exam ID", func(t *testing.T) {
		resp, err := handler.GetExamWithDetails(ctx, &pb.GetExamWithDetailsRequest{Id: "invalid"})
		assert.Error(t, err, "expected error for invalid exam ID")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "error code should be InvalidArgument")
	})

	t.Run("exam not found", func(t *testing.T) {
		examID := primitive.NewObjectID()
		examUC.EXPECT().GetExamWithDetails(ctx, examID).Return(nil, errors.New("not found"))

		resp, err := handler.GetExamWithDetails(ctx, &pb.GetExamWithDetailsRequest{Id: examID.Hex()})
		assert.Error(t, err, "expected error when exam not found")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.NotFound, st.Code(), "error code should be NotFound")
	})
}

func TestExamHandler_GenerateExamUsingAI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)
	handler := NewExamHandler(taskUC, questionUC, examUC)
	ctx := context.Background()

	t.Run("invalid user ID", func(t *testing.T) {
		resp, err := handler.GenerateExamUsingAI(ctx, &pb.GenerateExamUsingAIRequest{UserId: "invalid"})
		assert.Error(t, err, "expected error for invalid user ID")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "error code should be InvalidArgument")
	})

	t.Run("use case error", func(t *testing.T) {
		userID := primitive.NewObjectID()
		examUC.EXPECT().GenerateExamUsingAI(ctx, userID, 1, 1, "math", "10").Return(nil, errors.New("AI failure"))

		resp, err := handler.GenerateExamUsingAI(ctx, &pb.GenerateExamUsingAIRequest{
			UserId:       userID.Hex(),
			NumQuestions: 1,
			NumTasks:     1,
			Topic:        "math",
			Grade:        "10",
		})
		assert.Error(t, err, "expected error from use case")
		assert.Nil(t, resp, "response should be nil on error")
		st, ok := status.FromError(err)
		require.True(t, ok, "error should be gRPC status")
		assert.Equal(t, codes.Internal, st.Code(), "error code should be Internal")
	})
}
