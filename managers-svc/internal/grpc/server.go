package grpc

import (
	context "context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/core/domain"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/core/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	UnimplementedManagersServer
	service ports.ManagersService
	logger  *slog.Logger
}

func NewManagersServer(service ports.ManagersService, logger *slog.Logger) *Server {
	return &Server{
		service: service,
		logger:  logger,
	}
}

func StartGRPCServer(grpcPort string, managersService ports.ManagersService, logger *slog.Logger) error {
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	managersServer := NewManagersServer(managersService, logger)
	RegisterManagersServer(grpcServer, managersServer)
	reflection.Register(grpcServer)

	logger.Info("gRPC server listening", "port", grpcPort)
	return grpcServer.Serve(lis)
}

func NewServer(service ports.ManagersService) *Server {
	return &Server{service: service}
}

func (s *Server) AddGroupMember(ctx context.Context, req *AddMemberReq) (*GroupMemberResp, error) {
	s.logger.Info("Received AddGroupMember gRPC request",
		slog.String("group_id", req.GetGroupId()),
		slog.String("user_id", req.GetUserId()))

	groupID, err := uuid.Parse(req.GetGroupId())
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, err
	}

	role, err := domain.NewMemberRoleFromProto(req.GetRole().String())
	if err != nil {
		return nil, err
	}

	resp, err := s.service.AddGroupMember(ctx, groupID, userID, role)
	if err != nil {
		return nil, err
	}

	return &GroupMemberResp{
		Id:        resp.ID.String(),
		GroupId:   resp.GroupID.String(),
		UserId:    resp.UserID.String(),
		Role:      req.Role,
		CreatedAt: timestamppb.New(resp.CreatedAt),
		UpdatedAt: timestamppb.New(resp.UpdatedAt),
	}, nil
}

func (s *Server) AssignCourseInstructor(ctx context.Context, req *AssignInstructorReq) (*CourseInstructorResp, error) {
	s.logger.Info("Received AddGroupMember gRPC request",
		slog.String("course_id", req.GetCourseId()),
		slog.String("user_id", req.GetUserId()))

	courseID, err := uuid.Parse(req.GetCourseId())
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, err
	}
	resp, err := s.service.AssignCourseInstructor(ctx, courseID, userID)
	if err != nil {
		return nil, err
	}
	return &CourseInstructorResp{
		Id:        resp.ID.String(),
		CourseId:  resp.CourseID.String(),
		UserId:    resp.UserID.String(),
		CreatedAt: timestamppb.New(resp.CreatedAt),
		UpdatedAt: timestamppb.New(resp.UpdatedAt),
	}, nil
}
func (s *Server) CreateCourse(ctx context.Context, req *CreateCourseReq) (*CourseResp, error) {
	s.logger.Info("Received CreateCourse gRPC request",
		slog.String("name", req.GetName()))

	resp, err := s.service.CreateCourse(ctx, req.GetName())
	if err != nil {
		return nil, err
	}
	return &CourseResp{
		Id:        resp.ID.String(),
		Name:      resp.Name,
		CreatedAt: timestamppb.New(resp.CreatedAt),
		UpdatedAt: timestamppb.New(resp.UpdatedAt),
	}, nil
}
func (s *Server) CreateGroup(ctx context.Context, req *CreateGroupReq) (*GroupResp, error) {
	s.logger.Info("Received CreateGroup gRPC request",
		slog.String("course_id", req.GetCourseId()),
		slog.String("name", req.GetName()))

	courseID, err := uuid.Parse(req.GetCourseId())
	if err != nil {
		return nil, err
	}
	resp, err := s.service.CreateGroup(ctx, courseID, req.GetName())
	if err != nil {
		return nil, err
	}
	return &GroupResp{
		Id:        resp.ID.String(),
		CourseId:  resp.CourseID.String(),
		Name:      resp.Name,
		ExpireAt:  timestamppb.New(resp.ExpireAt),
		CreatedAt: timestamppb.New(resp.CreatedAt),
		UpdatedAt: timestamppb.New(resp.UpdatedAt),
	}, nil
}
func (s *Server) DeleteCourse(ctx context.Context, req *GetByIdReq) (*emptypb.Empty, error) {
	s.logger.Info("Received DeleteCourse gRPC request",
		slog.String("id", req.GetId()))

	ID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}
	err = s.service.DeleteCourse(ctx, ID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func (s *Server) DeleteGroup(ctx context.Context, req *GetByIdReq) (*emptypb.Empty, error) {
	s.logger.Info("Received DeleteGroup gRPC request",
		slog.String("id", req.GetId()))

	ID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}
	err = s.service.DeleteGroup(ctx, ID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func (s *Server) GetCourse(ctx context.Context, req *GetByIdReq) (*CourseResp, error) {
	s.logger.Info("Received GetCourse gRPC request",
		slog.String("id", req.GetId()))

	ID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}
	resp, err := s.service.GetCourse(ctx, ID)
	if err != nil {
		return nil, err
	}
	return &CourseResp{
		Id:        resp.ID.String(),
		Name:      resp.Name,
		CreatedAt: timestamppb.New(resp.CreatedAt),
		UpdatedAt: timestamppb.New(resp.UpdatedAt),
	}, nil
}
func (s *Server) GetGroup(ctx context.Context, req *GetByIdReq) (*GroupResp, error) {
	s.logger.Info("Received GetGroup gRPC request",
		slog.String("id", req.GetId()))

	ID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}
	resp, err := s.service.GetGroup(ctx, ID)
	if err != nil {
		return nil, err
	}
	return &GroupResp{
		Id:        resp.ID.String(),
		CourseId:  resp.CourseID.String(),
		Name:      resp.Name,
		ExpireAt:  timestamppb.New(resp.ExpireAt),
		CreatedAt: timestamppb.New(resp.CreatedAt),
		UpdatedAt: timestamppb.New(resp.UpdatedAt),
	}, nil
}
func (s *Server) IsUserGroupTeacher(ctx context.Context, req *CheckMemberReq) (*BoolResp, error) {
	s.logger.Info("Received IsUserGroupTeacher gRPC request",
		slog.String("group_id", req.GetGroupId()),
		slog.String("user_id", req.GetUserId()),
		slog.String("role", req.GetRole().String()))

	groupID, err := uuid.Parse(req.GetGroupId())
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, err
	}
	resp, err := s.service.IsUserGroupTeacher(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}
	return &BoolResp{
		Ok: resp,
	}, nil
}
func (s *Server) ListCourseInstructors(ctx context.Context, req *ListByCourseReq) (*ListCourseInstructorsResp, error) {
	s.logger.Info("Received ListCourseInstructors gRPC request",
		slog.String("course_id", req.GetCourseId()))

	courseID, err := uuid.Parse(req.GetCourseId())
	if err != nil {
		return nil, err
	}
	resp, err := s.service.ListCourseInstructors(ctx, courseID)
	if err != nil {
		return nil, err
	}
	instructors := make([]*CourseInstructorResp, len(resp))
	for i, inst := range resp {
		instructors[i] = &CourseInstructorResp{
			Id:        inst.ID.String(),
			CourseId:  req.GetCourseId(),
			UserId:    inst.UserID.String(),
			CreatedAt: timestamppb.New(inst.CreatedAt),
			UpdatedAt: timestamppb.New(inst.UpdatedAt),
		}
	}
	return &ListCourseInstructorsResp{
		Instructors: instructors,
	}, nil
}
func (s *Server) ListCourses(ctx context.Context, req *emptypb.Empty) (*ListCoursesResp, error) {
	s.logger.Info("Received ListCourses gRPC request")

	resp, err := s.service.ListCourses(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*CourseResp, len(resp))
	for i, item := range resp {
		items[i] = &CourseResp{
			Id:        item.ID.String(),
			Name:      item.Name,
			CreatedAt: timestamppb.New(item.CreatedAt),
			UpdatedAt: timestamppb.New(item.UpdatedAt),
		}
	}
	return &ListCoursesResp{
		Courses: items,
	}, nil
}
func (s *Server) ListGroupMembers(ctx context.Context, req *ListByGroupReq) (*ListGroupMembersResp, error) {
	s.logger.Info("Received ListGroupMembers gRPC request",
		slog.String("group_id", req.GetGroupId()))

	groupID, err := uuid.Parse(req.GetGroupId())
	if err != nil {
		return nil, err
	}
	resp, err := s.service.ListGroupMembers(ctx, groupID)
	if err != nil {
		return nil, err
	}
	members := make([]*GroupMemberResp, len(resp))
	for i, member := range resp {
		role, err := parseMemberRole(member.Role)
		if err != nil {
			return nil, err
		}
		members[i] = &GroupMemberResp{
			Id:        member.ID.String(),
			GroupId:   req.GetGroupId(),
			UserId:    member.UserID.String(),
			Role:      role,
			CreatedAt: timestamppb.New(member.CreatedAt),
			UpdatedAt: timestamppb.New(member.UpdatedAt),
		}
	}
	return &ListGroupMembersResp{
		Members: members,
	}, nil
}

func parseMemberRole(role domain.MemberRole) (MemberRole, error) {
	switch role {
	case domain.StudentRole:
		return MemberRole_STUDENT, nil
	case domain.TeacherRole:
		return MemberRole_TEACHER, nil
	default:
		return 0, fmt.Errorf("invalid role: %s", role)
	}
}

func (s *Server) ListGroupsByCourse(ctx context.Context, req *ListByCourseReq) (*ListGroupsResp, error) {
	s.logger.Info("Received ListGroupsByCourse gRPC request",
		slog.String("course_id", req.GetCourseId()))

	courseID, err := uuid.Parse(req.GetCourseId())
	if err != nil {
		return nil, err
	}
	resp, err := s.service.ListGroupsByCourse(ctx, courseID)
	if err != nil {
		return nil, err
	}
	groups := make([]*GroupResp, len(resp))
	for i, group := range resp {
		groups[i] = &GroupResp{
			Id:        group.ID.String(),
			CourseId:  req.GetCourseId(),
			Name:      group.Name,
			ExpireAt:  timestamppb.New(group.ExpireAt),
			CreatedAt: timestamppb.New(group.CreatedAt),
			UpdatedAt: timestamppb.New(group.UpdatedAt),
		}
	}
	return &ListGroupsResp{
		Groups: groups,
	}, nil
}

func (s *Server) ListStudentsByGroup(ctx context.Context, req *ListByGroupReq) (*ListGroupMembersResp, error) {
	return s.ListGroupMembers(ctx, req)
}

func (s *Server) ListTeachersByCourse(ctx context.Context, req *ListByCourseReq) (*ListCourseInstructorsResp, error) {
	return s.ListCourseInstructors(ctx, req)
}

func (s *Server) RemoveCourseInstructor(ctx context.Context, req *RemoveInstructorReq) (*emptypb.Empty, error) {
	s.logger.Info("Received RemoveCourseInstructor gRPC request",
		slog.String("course_id", req.GetCourseId()),
		slog.String("user_id", req.GetUserId()))

	courseID, err := uuid.Parse(req.GetCourseId())
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, err
	}
	if err := s.service.RemoveCourseInstructor(ctx, courseID, userID); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) RemoveGroupMember(ctx context.Context, req *RemoveMemberReq) (*emptypb.Empty, error) {
	s.logger.Info("Received RemoveGroupMember gRPC request",
		slog.String("group_id", req.GetGroupId()),
		slog.String("user_id", req.GetUserId()),
		slog.String("role", req.GetRole().String()))

	groupID, err := uuid.Parse(req.GetGroupId())
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, err
	}
	if err := s.service.RemoveGroupMember(ctx, groupID, userID, domain.MemberRole(req.GetRole().String())); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateCourse(ctx context.Context, req *UpdateCourseReq) (*CourseResp, error) {
	s.logger.Info("Received UpdateCourse gRPC request",
		slog.String("course_id", req.GetId()))

	courseID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}
	updatedCourse, err := s.service.UpdateCourse(ctx, domain.CourseUpdate{
		ID:   courseID,
		Name: &req.Name,
	})
	if err != nil {
		return nil, err
	}
	return &CourseResp{
		Id:        updatedCourse.ID.String(),
		Name:      updatedCourse.Name,
		CreatedAt: timestamppb.New(updatedCourse.CreatedAt),
		UpdatedAt: timestamppb.New(updatedCourse.UpdatedAt),
	}, nil
}

func (s *Server) UpdateGroup(ctx context.Context, req *UpdateGroupReq) (*GroupResp, error) {
	s.logger.Info("Received UpdateGroup gRPC request",
		slog.String("group_id", req.GetId()))

	groupID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	var expireAt time.Time
	if req.ExpireAt != nil {
		expireAt = req.ExpireAt.AsTime()
	}
	updatedGroup, err := s.service.UpdateGroup(ctx, domain.GroupUpdate{
		ID:       groupID,
		Name:     &req.Name,
		ExpireAt: &expireAt,
	})
	if err != nil {
		return nil, err
	}
	return &GroupResp{
		Id:        updatedGroup.ID.String(),
		CourseId:  updatedGroup.CourseID.String(),
		Name:      updatedGroup.Name,
		ExpireAt:  timestamppb.New(updatedGroup.ExpireAt),
		CreatedAt: timestamppb.New(updatedGroup.CreatedAt),
		UpdatedAt: timestamppb.New(updatedGroup.UpdatedAt),
	}, nil
}
