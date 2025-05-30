// Code generated by MockGen. DO NOT EDIT.
// Source: repository/interface.go
//
// Generated by this command:
//
//	mockgen -source=repository/interface.go -destination=repository/mocks/repository_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	gomock "go.uber.org/mock/gomock"
)

// MockTaskRepository is a mock of TaskRepository interface.
type MockTaskRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTaskRepositoryMockRecorder
	isgomock struct{}
}

// MockTaskRepositoryMockRecorder is the mock recorder for MockTaskRepository.
type MockTaskRepositoryMockRecorder struct {
	mock *MockTaskRepository
}

// NewMockTaskRepository creates a new mock instance.
func NewMockTaskRepository(ctrl *gomock.Controller) *MockTaskRepository {
	mock := &MockTaskRepository{ctrl: ctrl}
	mock.recorder = &MockTaskRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskRepository) EXPECT() *MockTaskRepositoryMockRecorder {
	return m.recorder
}

// CreateTask mocks base method.
func (m *MockTaskRepository) CreateTask(ctx context.Context, task *domain.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", ctx, task)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockTaskRepositoryMockRecorder) CreateTask(ctx, task any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockTaskRepository)(nil).CreateTask), ctx, task)
}

// CreateTaskWithTransaction mocks base method.
func (m *MockTaskRepository) CreateTaskWithTransaction(ctx context.Context, task *domain.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTaskWithTransaction", ctx, task)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTaskWithTransaction indicates an expected call of CreateTaskWithTransaction.
func (mr *MockTaskRepositoryMockRecorder) CreateTaskWithTransaction(ctx, task any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTaskWithTransaction", reflect.TypeOf((*MockTaskRepository)(nil).CreateTaskWithTransaction), ctx, task)
}

// DeleteTask mocks base method.
func (m *MockTaskRepository) DeleteTask(ctx context.Context, id primitive.ObjectID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTask", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTask indicates an expected call of DeleteTask.
func (mr *MockTaskRepositoryMockRecorder) DeleteTask(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockTaskRepository)(nil).DeleteTask), ctx, id)
}

// GetAllTasks mocks base method.
func (m *MockTaskRepository) GetAllTasks(ctx context.Context) ([]domain.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTasks", ctx)
	ret0, _ := ret[0].([]domain.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllTasks indicates an expected call of GetAllTasks.
func (mr *MockTaskRepositoryMockRecorder) GetAllTasks(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTasks", reflect.TypeOf((*MockTaskRepository)(nil).GetAllTasks), ctx)
}

// GetTaskByID mocks base method.
func (m *MockTaskRepository) GetTaskByID(ctx context.Context, id primitive.ObjectID) (*domain.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTaskByID", ctx, id)
	ret0, _ := ret[0].(*domain.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskByID indicates an expected call of GetTaskByID.
func (mr *MockTaskRepositoryMockRecorder) GetTaskByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskByID", reflect.TypeOf((*MockTaskRepository)(nil).GetTaskByID), ctx, id)
}

// GetTasksByExamID mocks base method.
func (m *MockTaskRepository) GetTasksByExamID(ctx context.Context, examID primitive.ObjectID) ([]domain.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTasksByExamID", ctx, examID)
	ret0, _ := ret[0].([]domain.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTasksByExamID indicates an expected call of GetTasksByExamID.
func (mr *MockTaskRepositoryMockRecorder) GetTasksByExamID(ctx, examID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTasksByExamID", reflect.TypeOf((*MockTaskRepository)(nil).GetTasksByExamID), ctx, examID)
}

// UpdateTask mocks base method.
func (m *MockTaskRepository) UpdateTask(ctx context.Context, task *domain.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", ctx, task)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockTaskRepositoryMockRecorder) UpdateTask(ctx, task any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockTaskRepository)(nil).UpdateTask), ctx, task)
}

// MockQuestionRepository is a mock of QuestionRepository interface.
type MockQuestionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockQuestionRepositoryMockRecorder
	isgomock struct{}
}

// MockQuestionRepositoryMockRecorder is the mock recorder for MockQuestionRepository.
type MockQuestionRepositoryMockRecorder struct {
	mock *MockQuestionRepository
}

// NewMockQuestionRepository creates a new mock instance.
func NewMockQuestionRepository(ctrl *gomock.Controller) *MockQuestionRepository {
	mock := &MockQuestionRepository{ctrl: ctrl}
	mock.recorder = &MockQuestionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuestionRepository) EXPECT() *MockQuestionRepositoryMockRecorder {
	return m.recorder
}

// CreateQuestion mocks base method.
func (m *MockQuestionRepository) CreateQuestion(ctx context.Context, question *domain.Question) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateQuestion", ctx, question)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateQuestion indicates an expected call of CreateQuestion.
func (mr *MockQuestionRepositoryMockRecorder) CreateQuestion(ctx, question any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateQuestion", reflect.TypeOf((*MockQuestionRepository)(nil).CreateQuestion), ctx, question)
}

// CreateQuestionWithTransaction mocks base method.
func (m *MockQuestionRepository) CreateQuestionWithTransaction(ctx context.Context, question *domain.Question) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateQuestionWithTransaction", ctx, question)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateQuestionWithTransaction indicates an expected call of CreateQuestionWithTransaction.
func (mr *MockQuestionRepositoryMockRecorder) CreateQuestionWithTransaction(ctx, question any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateQuestionWithTransaction", reflect.TypeOf((*MockQuestionRepository)(nil).CreateQuestionWithTransaction), ctx, question)
}

// DeleteQuestion mocks base method.
func (m *MockQuestionRepository) DeleteQuestion(ctx context.Context, id primitive.ObjectID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteQuestion", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteQuestion indicates an expected call of DeleteQuestion.
func (mr *MockQuestionRepositoryMockRecorder) DeleteQuestion(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteQuestion", reflect.TypeOf((*MockQuestionRepository)(nil).DeleteQuestion), ctx, id)
}

// GetAllQuestions mocks base method.
func (m *MockQuestionRepository) GetAllQuestions(ctx context.Context) ([]domain.Question, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllQuestions", ctx)
	ret0, _ := ret[0].([]domain.Question)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllQuestions indicates an expected call of GetAllQuestions.
func (mr *MockQuestionRepositoryMockRecorder) GetAllQuestions(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllQuestions", reflect.TypeOf((*MockQuestionRepository)(nil).GetAllQuestions), ctx)
}

// GetQuestionByID mocks base method.
func (m *MockQuestionRepository) GetQuestionByID(ctx context.Context, id primitive.ObjectID) (*domain.Question, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuestionByID", ctx, id)
	ret0, _ := ret[0].(*domain.Question)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuestionByID indicates an expected call of GetQuestionByID.
func (mr *MockQuestionRepositoryMockRecorder) GetQuestionByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuestionByID", reflect.TypeOf((*MockQuestionRepository)(nil).GetQuestionByID), ctx, id)
}

// GetQuestionsByExamID mocks base method.
func (m *MockQuestionRepository) GetQuestionsByExamID(ctx context.Context, examID primitive.ObjectID) ([]domain.Question, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuestionsByExamID", ctx, examID)
	ret0, _ := ret[0].([]domain.Question)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuestionsByExamID indicates an expected call of GetQuestionsByExamID.
func (mr *MockQuestionRepositoryMockRecorder) GetQuestionsByExamID(ctx, examID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuestionsByExamID", reflect.TypeOf((*MockQuestionRepository)(nil).GetQuestionsByExamID), ctx, examID)
}

// UpdateQuestion mocks base method.
func (m *MockQuestionRepository) UpdateQuestion(ctx context.Context, question *domain.Question) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateQuestion", ctx, question)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateQuestion indicates an expected call of UpdateQuestion.
func (mr *MockQuestionRepositoryMockRecorder) UpdateQuestion(ctx, question any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateQuestion", reflect.TypeOf((*MockQuestionRepository)(nil).UpdateQuestion), ctx, question)
}

// MockExamRepository is a mock of ExamRepository interface.
type MockExamRepository struct {
	ctrl     *gomock.Controller
	recorder *MockExamRepositoryMockRecorder
	isgomock struct{}
}

// MockExamRepositoryMockRecorder is the mock recorder for MockExamRepository.
type MockExamRepositoryMockRecorder struct {
	mock *MockExamRepository
}

// NewMockExamRepository creates a new mock instance.
func NewMockExamRepository(ctrl *gomock.Controller) *MockExamRepository {
	mock := &MockExamRepository{ctrl: ctrl}
	mock.recorder = &MockExamRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExamRepository) EXPECT() *MockExamRepositoryMockRecorder {
	return m.recorder
}

// CreateExam mocks base method.
func (m *MockExamRepository) CreateExam(ctx context.Context, exam *domain.Exam) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateExam", ctx, exam)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateExam indicates an expected call of CreateExam.
func (mr *MockExamRepositoryMockRecorder) CreateExam(ctx, exam any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateExam", reflect.TypeOf((*MockExamRepository)(nil).CreateExam), ctx, exam)
}

// DeleteExam mocks base method.
func (m *MockExamRepository) DeleteExam(ctx context.Context, id primitive.ObjectID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteExam", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteExam indicates an expected call of DeleteExam.
func (mr *MockExamRepositoryMockRecorder) DeleteExam(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteExam", reflect.TypeOf((*MockExamRepository)(nil).DeleteExam), ctx, id)
}

// DeleteExamWithTransaction mocks base method.
func (m *MockExamRepository) DeleteExamWithTransaction(ctx context.Context, id primitive.ObjectID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteExamWithTransaction", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteExamWithTransaction indicates an expected call of DeleteExamWithTransaction.
func (mr *MockExamRepositoryMockRecorder) DeleteExamWithTransaction(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteExamWithTransaction", reflect.TypeOf((*MockExamRepository)(nil).DeleteExamWithTransaction), ctx, id)
}

// GetAllExams mocks base method.
func (m *MockExamRepository) GetAllExams(ctx context.Context) ([]domain.Exam, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllExams", ctx)
	ret0, _ := ret[0].([]domain.Exam)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllExams indicates an expected call of GetAllExams.
func (mr *MockExamRepositoryMockRecorder) GetAllExams(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllExams", reflect.TypeOf((*MockExamRepository)(nil).GetAllExams), ctx)
}

// GetExamByID mocks base method.
func (m *MockExamRepository) GetExamByID(ctx context.Context, id primitive.ObjectID) (*domain.Exam, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExamByID", ctx, id)
	ret0, _ := ret[0].(*domain.Exam)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExamByID indicates an expected call of GetExamByID.
func (mr *MockExamRepositoryMockRecorder) GetExamByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExamByID", reflect.TypeOf((*MockExamRepository)(nil).GetExamByID), ctx, id)
}

// GetExamsByUser mocks base method.
func (m *MockExamRepository) GetExamsByUser(ctx context.Context, userID primitive.ObjectID) ([]domain.Exam, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExamsByUser", ctx, userID)
	ret0, _ := ret[0].([]domain.Exam)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExamsByUser indicates an expected call of GetExamsByUser.
func (mr *MockExamRepositoryMockRecorder) GetExamsByUser(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExamsByUser", reflect.TypeOf((*MockExamRepository)(nil).GetExamsByUser), ctx, userID)
}

// UpdateExam mocks base method.
func (m *MockExamRepository) UpdateExam(ctx context.Context, exam *domain.Exam) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateExam", ctx, exam)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateExam indicates an expected call of UpdateExam.
func (mr *MockExamRepositoryMockRecorder) UpdateExam(ctx, exam any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateExam", reflect.TypeOf((*MockExamRepository)(nil).UpdateExam), ctx, exam)
}

// UpdateExamStatus mocks base method.
func (m *MockExamRepository) UpdateExamStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateExamStatus", ctx, id, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateExamStatus indicates an expected call of UpdateExamStatus.
func (mr *MockExamRepositoryMockRecorder) UpdateExamStatus(ctx, id, status any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateExamStatus", reflect.TypeOf((*MockExamRepository)(nil).UpdateExamStatus), ctx, id, status)
}
