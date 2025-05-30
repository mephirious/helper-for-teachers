package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/core/domain"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/core/ports"
)

type ManagersService struct {
	timeSource            func() time.Time
	logger                *slog.Logger
	courseRepository      ports.CourseRepository
	groupRepository       ports.GroupRepository
	memberRepository      ports.MemberRepository
	instructorsRepository ports.InstructorRepository
	publisher             ports.NatsPublisher
	SMTP                  ports.SMTP
}

func NewManagersService(timeSource func() time.Time, logger *slog.Logger, courseRepo ports.CourseRepository, groupRepository ports.GroupRepository, memberRepository ports.MemberRepository, instructorsRepository ports.InstructorRepository, publisher ports.NatsPublisher, smtp ports.SMTP) *ManagersService {
	return &ManagersService{
		timeSource:            timeSource,
		logger:                logger,
		courseRepository:      courseRepo,
		groupRepository:       groupRepository,
		memberRepository:      memberRepository,
		instructorsRepository: instructorsRepository,
		publisher:             publisher,
		SMTP:                  smtp,
	}
}

func validateCourseName(name string) error {
	if len(name) < 3 {
		return fmt.Errorf("course name must be at least 3 characters")
	}
	return nil
}

func validateGroupName(name string) error {
	if len(name) < 3 {
		return fmt.Errorf("group name must be at least 3 characters")
	}
	return nil
}

func (s *ManagersService) CreateCourse(ctx context.Context, name string) (*domain.Course, error) {
	const op = "ManagersService.CreateCourse"
	if err := validateCourseName(name); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	course := &domain.Course{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: s.timeSource(),
		UpdatedAt: s.timeSource(),
	}
	if err := s.courseRepository.Create(ctx, course); err != nil {
		return nil, fmt.Errorf("CreateCourse: %w", err)
	}
	if s.publisher != nil {
		s.publisher.PublishCourseCreated(ctx, course)
	}
	return course, nil
}

func (s *ManagersService) GetCourse(ctx context.Context, id uuid.UUID) (*domain.Course, error) {
	course, err := s.courseRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetCourse: %w", err)
	}
	return course, nil
}

func (s *ManagersService) UpdateCourse(ctx context.Context, req domain.CourseUpdate) (*domain.Course, error) {
	const op = "ManagersService.UpdateCourse"
	var courseData *domain.Course
	if req.Name != nil && *req.Name != "" {
		if err := validateCourseName(*req.Name); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	return courseData, s.courseRepository.Update(ctx, req.ID, func(item *domain.Course) (updated bool, err error) {
		if req.Name != nil && *req.Name != item.Name {
			item.Name = *req.Name
			updated = true
		}

		if !updated {
			return
		}

		item.UpdatedAt = s.timeSource().UTC()
		courseData = item

		return
	})
}

func (s *ManagersService) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	return s.courseRepository.Delete(ctx, id)
}

func (s *ManagersService) ListCourses(ctx context.Context) ([]*domain.Course, error) {
	return s.courseRepository.List(ctx)
}

func (s *ManagersService) CreateGroup(ctx context.Context, courseID uuid.UUID, name string) (*domain.Group, error) {
	const op = "ManagersService.CreateGroup"
	if err := validateGroupName(name); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	group := &domain.Group{
		ID:        uuid.New(),
		CourseID:  courseID,
		Name:      name,
		ExpireAt:  s.timeSource().Add(time.Hour * 24 * 7 * 10),
		CreatedAt: s.timeSource(),
		UpdatedAt: s.timeSource(),
	}

	if _, err := s.courseRepository.GetByID(ctx, courseID); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.groupRepository.Create(ctx, group); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if s.publisher != nil {
		s.publisher.PublishGroupCreated(ctx, group)
	}
	return group, nil
}

func (s *ManagersService) GetGroup(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
	return s.groupRepository.GetByID(ctx, id)
}

func (s *ManagersService) UpdateGroup(ctx context.Context, req domain.GroupUpdate) (*domain.Group, error) {
	const op = "ManagersService.UpdateCourse"
	var groupData *domain.Group

	return groupData, s.groupRepository.Update(ctx, req.ID, func(item *domain.Group) (updated bool, err error) {
		if req.Name != nil && *req.Name != item.Name {
			if err := validateCourseName(*req.Name); err != nil {
				return false, fmt.Errorf("%s: %w", op, err)
			}
			item.Name = *req.Name
			updated = true
		}

		if req.ExpireAt != nil && item.ExpireAt.Compare(*req.ExpireAt) != 0 && req.ExpireAt.After(s.timeSource()) {
			item.ExpireAt = *req.ExpireAt
			updated = true
		}

		if !updated {
			return
		}

		item.UpdatedAt = s.timeSource().UTC()
		groupData = item

		return
	})
}

func (s *ManagersService) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	return s.groupRepository.Delete(ctx, id)
}

func (s *ManagersService) ListGroupsByCourse(ctx context.Context, courseID uuid.UUID) ([]*domain.Group, error) {
	return s.groupRepository.ListByCourse(ctx, courseID)
}

func (s *ManagersService) AddGroupMember(ctx context.Context, groupID, userID uuid.UUID, role domain.MemberRole) (*domain.GroupMember, error) {
	const op = "ManagersService.AddGroupMember"
	member := &domain.GroupMember{
		ID:        uuid.New(),
		GroupID:   groupID,
		UserID:    userID,
		Role:      role,
		CreatedAt: s.timeSource(),
		UpdatedAt: s.timeSource(),
	}

	if _, err := s.groupRepository.GetByID(ctx, member.GroupID); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.memberRepository.Create(ctx, member); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if s.publisher != nil {
		s.publisher.PublishGroupMemberAdded(ctx, member)
	}
	return member, nil
}

func (s *ManagersService) RemoveGroupMember(ctx context.Context, groupID, userID uuid.UUID, role domain.MemberRole) error {
	return s.memberRepository.Delete(ctx, groupID, userID, role)
}

func (s *ManagersService) ListGroupMembers(ctx context.Context, groupID uuid.UUID) ([]*domain.GroupMember, error) {
	return s.memberRepository.ListByGroup(ctx, groupID)
}

func (s *ManagersService) IsUserGroupTeacher(ctx context.Context, groupID, userID uuid.UUID) (bool, error) {
	return s.memberRepository.Exists(ctx, groupID, userID, domain.TeacherRole)
}

func (s *ManagersService) AssignCourseInstructor(ctx context.Context, courseID, userID uuid.UUID) (*domain.CourseInstructor, error) {
	const op = "ManagersService.AssignCourseInstructor"
	instructor := &domain.CourseInstructor{
		ID:        uuid.New(),
		CourseID:  courseID,
		UserID:    userID,
		CreatedAt: s.timeSource(),
		UpdatedAt: s.timeSource(),
	}
	if _, err := s.courseRepository.GetByID(ctx, courseID); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.instructorsRepository.Create(ctx, instructor); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if s.publisher != nil {
		s.publisher.PublishCourseInstructorAssigned(ctx, instructor)
	}
	return instructor, nil
}

func (s *ManagersService) RemoveCourseInstructor(ctx context.Context, courseID, userID uuid.UUID) error {
	return s.instructorsRepository.Delete(ctx, courseID, userID)
}

func (s *ManagersService) ListCourseInstructors(ctx context.Context, courseID uuid.UUID) ([]*domain.CourseInstructor, error) {
	return s.instructorsRepository.ListByCourse(ctx, courseID)
}
