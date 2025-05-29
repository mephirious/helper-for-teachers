package handler

import (
	"context"
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/usecase"
	pb "github.com/mephirious/helper-for-teachers/services/exam-svc/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ExamHandler struct {
	pb.UnimplementedExamServiceServer
	taskUseCase     usecase.TaskUseCase
	questionUseCase usecase.QuestionUseCase
	examUseCase     usecase.ExamUseCase
}

func NewExamHandler(taskUseCase usecase.TaskUseCase, questionUseCase usecase.QuestionUseCase, examUseCase usecase.ExamUseCase) *ExamHandler {
	return &ExamHandler{
		taskUseCase:     taskUseCase,
		questionUseCase: questionUseCase,
		examUseCase:     examUseCase,
	}
}

func (h *ExamHandler) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.TaskResponse, error) {
	examID, err := primitive.ObjectIDFromHex(req.Task.ExamId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid exam_id: %v", err)
	}

	task := &domain.Task{
		ExamID:      examID,
		TaskType:    req.Task.TaskType,
		Description: req.Task.Description,
		Score:       req.Task.Score,
		CreatedAt:   time.Now(),
	}

	createdTask, err := h.taskUseCase.CreateTask(ctx, task)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create task: %v", err)
	}

	return &pb.TaskResponse{
		Id:          createdTask.ID.Hex(),
		ExamId:      createdTask.ExamID.Hex(),
		TaskType:    createdTask.TaskType,
		Description: createdTask.Description,
		Score:       createdTask.Score,
		CreatedAt:   timestamppb.New(createdTask.CreatedAt),
	}, nil
}

func (h *ExamHandler) GetTaskByID(ctx context.Context, req *pb.GetTaskByIDRequest) (*pb.TaskResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	task, err := h.taskUseCase.GetTaskByID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "task not found: %v", err)
	}

	return &pb.TaskResponse{
		Id:          task.ID.Hex(),
		ExamId:      task.ExamID.Hex(),
		TaskType:    task.TaskType,
		Description: task.Description,
		Score:       task.Score,
		CreatedAt:   timestamppb.New(task.CreatedAt),
	}, nil
}

func (h *ExamHandler) GetTasksByExamID(ctx context.Context, req *pb.GetTasksByExamIDRequest) (*pb.TasksResponse, error) {
	examID, err := primitive.ObjectIDFromHex(req.ExamId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid exam_id: %v", err)
	}

	tasks, err := h.taskUseCase.GetTasksByExamID(ctx, examID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get tasks: %v", err)
	}

	protoTasks := make([]*pb.Task, len(tasks))
	for i, task := range tasks {
		protoTasks[i] = &pb.Task{
			Id:          task.ID.Hex(),
			ExamId:      task.ExamID.Hex(),
			TaskType:    task.TaskType,
			Description: task.Description,
			Score:       task.Score,
			CreatedAt:   timestamppb.New(task.CreatedAt),
		}
	}

	return &pb.TasksResponse{Tasks: protoTasks}, nil
}

func (h *ExamHandler) GetAllTasks(ctx context.Context, _ *emptypb.Empty) (*pb.TasksResponse, error) {
	tasks, err := h.taskUseCase.GetAllTasks(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get all tasks: %v", err)
	}

	protoTasks := make([]*pb.Task, len(tasks))
	for i, task := range tasks {
		protoTasks[i] = &pb.Task{
			Id:          task.ID.Hex(),
			ExamId:      task.ExamID.Hex(),
			TaskType:    task.TaskType,
			Description: task.Description,
			Score:       task.Score,
			CreatedAt:   timestamppb.New(task.CreatedAt),
		}
	}

	return &pb.TasksResponse{Tasks: protoTasks}, nil
}

func (h *ExamHandler) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*emptypb.Empty, error) {
	id, err := primitive.ObjectIDFromHex(req.Task.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	examID, err := primitive.ObjectIDFromHex(req.Task.ExamId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid exam_id: %v", err)
	}

	task := &domain.Task{
		ID:          id,
		ExamID:      examID,
		TaskType:    req.Task.TaskType,
		Description: req.Task.Description,
		Score:       req.Task.Score,
		CreatedAt:   req.Task.CreatedAt.AsTime(),
	}

	if err := h.taskUseCase.UpdateTask(ctx, task); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update task: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *ExamHandler) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*emptypb.Empty, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	if err := h.taskUseCase.DeleteTask(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete task: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *ExamHandler) CreateQuestion(ctx context.Context, req *pb.CreateQuestionRequest) (*pb.QuestionResponse, error) {
	examID, err := primitive.ObjectIDFromHex(req.Question.ExamId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid exam_id: %v", err)
	}

	question := &domain.Question{
		ExamID:        examID,
		QuestionText:  req.Question.QuestionText,
		Options:       req.Question.Options,
		CorrectAnswer: req.Question.CorrectAnswer,
		Status:        req.Question.Status,
		CreatedAt:     time.Now(),
	}

	createdQuestion, err := h.questionUseCase.CreateQuestion(ctx, question)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create question: %v", err)
	}

	return &pb.QuestionResponse{
		Id:            createdQuestion.ID.Hex(),
		ExamId:        createdQuestion.ExamID.Hex(),
		QuestionText:  createdQuestion.QuestionText,
		Options:       createdQuestion.Options,
		CorrectAnswer: createdQuestion.CorrectAnswer,
		Status:        createdQuestion.Status,
		CreatedAt:     timestamppb.New(createdQuestion.CreatedAt),
	}, nil
}

func (h *ExamHandler) GetQuestionByID(ctx context.Context, req *pb.GetQuestionByIDRequest) (*pb.QuestionResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	question, err := h.questionUseCase.GetQuestionByID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "question not found: %v", err)
	}

	return &pb.QuestionResponse{
		Id:            question.ID.Hex(),
		ExamId:        question.ExamID.Hex(),
		QuestionText:  question.QuestionText,
		Options:       question.Options,
		CorrectAnswer: question.CorrectAnswer,
		Status:        question.Status,
		CreatedAt:     timestamppb.New(question.CreatedAt),
	}, nil
}

func (h *ExamHandler) GetQuestionsByExamID(ctx context.Context, req *pb.GetQuestionsByExamIDRequest) (*pb.QuestionsResponse, error) {
	examID, err := primitive.ObjectIDFromHex(req.ExamId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid exam_id: %v", err)
	}

	questions, err := h.questionUseCase.GetQuestionsByExamID(ctx, examID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get questions: %v", err)
	}

	protoQuestions := make([]*pb.Question, len(questions))
	for i, question := range questions {
		protoQuestions[i] = &pb.Question{
			Id:            question.ID.Hex(),
			ExamId:        question.ExamID.Hex(),
			QuestionText:  question.QuestionText,
			Options:       question.Options,
			CorrectAnswer: question.CorrectAnswer,
			Status:        question.Status,
			CreatedAt:     timestamppb.New(question.CreatedAt),
		}
	}

	return &pb.QuestionsResponse{Questions: protoQuestions}, nil
}

func (h *ExamHandler) GetAllQuestions(ctx context.Context, _ *emptypb.Empty) (*pb.QuestionsResponse, error) {
	questions, err := h.questionUseCase.GetAllQuestions(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get all questions: %v", err)
	}

	protoQuestions := make([]*pb.Question, len(questions))
	for i, question := range questions {
		protoQuestions[i] = &pb.Question{
			Id:            question.ID.Hex(),
			ExamId:        question.ExamID.Hex(),
			QuestionText:  question.QuestionText,
			Options:       question.Options,
			CorrectAnswer: question.CorrectAnswer,
			Status:        question.Status,
			CreatedAt:     timestamppb.New(question.CreatedAt),
		}
	}

	return &pb.QuestionsResponse{Questions: protoQuestions}, nil
}

func (h *ExamHandler) UpdateQuestion(ctx context.Context, req *pb.UpdateQuestionRequest) (*emptypb.Empty, error) {
	id, err := primitive.ObjectIDFromHex(req.Question.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	examID, err := primitive.ObjectIDFromHex(req.Question.ExamId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid exam_id: %v", err)
	}

	question := &domain.Question{
		ID:            id,
		ExamID:        examID,
		QuestionText:  req.Question.QuestionText,
		Options:       req.Question.Options,
		CorrectAnswer: req.Question.CorrectAnswer,
		Status:        req.Question.Status,
		CreatedAt:     req.Question.CreatedAt.AsTime(),
	}

	if err := h.questionUseCase.UpdateQuestion(ctx, question); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update question: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *ExamHandler) DeleteQuestion(ctx context.Context, req *pb.DeleteQuestionRequest) (*emptypb.Empty, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	if err := h.questionUseCase.DeleteQuestion(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete question: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *ExamHandler) CreateExam(ctx context.Context, req *pb.CreateExamRequest) (*pb.ExamResponse, error) {
	createdBy, err := primitive.ObjectIDFromHex(req.Exam.CreatedBy)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid created_by: %v", err)
	}

	exam := &domain.Exam{
		Title:       req.Exam.Title,
		Description: req.Exam.Description,
		CreatedBy:   createdBy,
		Status:      req.Exam.Status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdExam, err := h.examUseCase.CreateExam(ctx, exam)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create exam: %v", err)
	}

	return &pb.ExamResponse{
		Id:          createdExam.ID.Hex(),
		Title:       createdExam.Title,
		Description: createdExam.Description,
		CreatedBy:   createdExam.CreatedBy.Hex(),
		Status:      createdExam.Status,
		CreatedAt:   timestamppb.New(createdExam.CreatedAt),
		UpdatedAt:   timestamppb.New(createdExam.UpdatedAt),
	}, nil
}

func (h *ExamHandler) GetExamByID(ctx context.Context, req *pb.GetExamByIDRequest) (*pb.ExamResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	exam, err := h.examUseCase.GetExamByID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "exam not found: %v", err)
	}

	return &pb.ExamResponse{
		Id:          exam.ID.Hex(),
		Title:       exam.Title,
		Description: exam.Description,
		CreatedBy:   exam.CreatedBy.Hex(),
		Status:      exam.Status,
		CreatedAt:   timestamppb.New(exam.CreatedAt),
		UpdatedAt:   timestamppb.New(exam.UpdatedAt),
	}, nil
}

func (h *ExamHandler) GetExamsByUser(ctx context.Context, req *pb.GetExamsByUserRequest) (*pb.ExamsResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	exams, err := h.examUseCase.GetExamsByUser(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get exams: %v", err)
	}

	protoExams := make([]*pb.Exam, len(exams))
	for i, exam := range exams {
		protoExams[i] = &pb.Exam{
			Id:          exam.ID.Hex(),
			Title:       exam.Title,
			Description: exam.Description,
			CreatedBy:   exam.CreatedBy.Hex(),
			Status:      exam.Status,
			CreatedAt:   timestamppb.New(exam.CreatedAt),
			UpdatedAt:   timestamppb.New(exam.UpdatedAt),
		}
	}

	return &pb.ExamsResponse{Exams: protoExams}, nil
}

func (h *ExamHandler) UpdateExam(ctx context.Context, req *pb.UpdateExamRequest) (*emptypb.Empty, error) {
	id, err := primitive.ObjectIDFromHex(req.Exam.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	createdBy, err := primitive.ObjectIDFromHex(req.Exam.CreatedBy)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid created_by: %v", err)
	}

	exam := &domain.Exam{
		ID:          id,
		Title:       req.Exam.Title,
		Description: req.Exam.Description,
		CreatedBy:   createdBy,
		Status:      req.Exam.Status,
		CreatedAt:   req.Exam.CreatedAt.AsTime(),
		UpdatedAt:   time.Now(),
	}

	if err := h.examUseCase.UpdateExam(ctx, exam); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update exam: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *ExamHandler) UpdateExamStatus(ctx context.Context, req *pb.UpdateExamStatusRequest) (*emptypb.Empty, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	if err := h.examUseCase.UpdateExamStatus(ctx, id, req.Status); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update exam status: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *ExamHandler) DeleteExam(ctx context.Context, req *pb.DeleteExamRequest) (*emptypb.Empty, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	if err := h.examUseCase.DeleteExam(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete exam: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *ExamHandler) GetAllExams(ctx context.Context, _ *emptypb.Empty) (*pb.ExamsResponse, error) {
	exams, err := h.examUseCase.GetAllExams(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get all exams: %v", err)
	}

	protoExams := make([]*pb.Exam, len(exams))
	for i, exam := range exams {
		protoExams[i] = &pb.Exam{
			Id:          exam.ID.Hex(),
			Title:       exam.Title,
			Description: exam.Description,
			CreatedBy:   exam.CreatedBy.Hex(),
			Status:      exam.Status,
			CreatedAt:   timestamppb.New(exam.CreatedAt),
			UpdatedAt:   timestamppb.New(exam.UpdatedAt),
		}
	}

	return &pb.ExamsResponse{Exams: protoExams}, nil
}

func (h *ExamHandler) GetExamWithDetails(ctx context.Context, req *pb.GetExamWithDetailsRequest) (*pb.ExamDetailedResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	examDetailed, err := h.examUseCase.GetExamWithDetails(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "exam not found: %v", err)
	}

	protoTasks := make([]*pb.Task, len(examDetailed.Tasks))
	for i, task := range examDetailed.Tasks {
		protoTasks[i] = &pb.Task{
			Id:          task.ID.Hex(),
			ExamId:      task.ExamID.Hex(),
			TaskType:    task.TaskType,
			Description: task.Description,
			Score:       task.Score,
			CreatedAt:   timestamppb.New(task.CreatedAt),
		}
	}

	protoQuestions := make([]*pb.Question, len(examDetailed.Questions))
	for i, question := range examDetailed.Questions {
		protoQuestions[i] = &pb.Question{
			Id:            question.ID.Hex(),
			ExamId:        question.ExamID.Hex(),
			QuestionText:  question.QuestionText,
			Options:       question.Options,
			CorrectAnswer: question.CorrectAnswer,
			Status:        question.Status,
			CreatedAt:     timestamppb.New(question.CreatedAt),
		}
	}

	return &pb.ExamDetailedResponse{
		Id:          examDetailed.ID.Hex(),
		Title:       examDetailed.Title,
		Description: examDetailed.Description,
		CreatedBy:   examDetailed.CreatedBy.Hex(),
		Status:      examDetailed.Status,
		CreatedAt:   timestamppb.New(examDetailed.CreatedAt),
		UpdatedAt:   timestamppb.New(examDetailed.UpdatedAt),
		Tasks:       protoTasks,
		Questions:   protoQuestions,
	}, nil
}

func (h *ExamHandler) GenerateExamUsingAI(ctx context.Context, req *pb.GenerateExamUsingAIRequest) (*pb.ExamDetailedResponse, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	examDetailed, err := h.examUseCase.GenerateExamUsingAI(ctx, userID, int(req.NumQuestions), int(req.NumTasks), req.Topic, req.Grade)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate exam: %v", err)
	}

	protoTasks := make([]*pb.Task, len(examDetailed.Tasks))
	for i, task := range examDetailed.Tasks {
		protoTasks[i] = &pb.Task{
			Id:          task.ID.Hex(),
			ExamId:      task.ExamID.Hex(),
			TaskType:    task.TaskType,
			Description: task.Description,
			Score:       task.Score,
			CreatedAt:   timestamppb.New(task.CreatedAt),
		}
	}

	protoQuestions := make([]*pb.Question, len(examDetailed.Questions))
	for i, question := range examDetailed.Questions {
		protoQuestions[i] = &pb.Question{
			Id:            question.ID.Hex(),
			ExamId:        question.ExamID.Hex(),
			QuestionText:  question.QuestionText,
			Options:       question.Options,
			CorrectAnswer: question.CorrectAnswer,
			Status:        question.Status,
			CreatedAt:     timestamppb.New(question.CreatedAt),
		}
	}

	return &pb.ExamDetailedResponse{
		Id:          examDetailed.ID.Hex(),
		Title:       examDetailed.Title,
		Description: examDetailed.Description,
		CreatedBy:   examDetailed.CreatedBy.Hex(),
		Status:      examDetailed.Status,
		CreatedAt:   timestamppb.New(examDetailed.CreatedAt),
		UpdatedAt:   timestamppb.New(examDetailed.UpdatedAt),
		Tasks:       protoTasks,
		Questions:   protoQuestions,
	}, nil
}
