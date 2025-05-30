package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/core/domain"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/core/service"
)

type (
	MockCourseRepository struct {
		courses map[uuid.UUID]*domain.Course
		err     error
	}

	MockGroupRepository struct {
		groups map[uuid.UUID]*domain.Group
		err    error
	}

	MockMemberRepository struct {
		members []*domain.GroupMember
		exists  bool
		err     error
	}

	MockInstructorRepository struct {
		instructors []*domain.CourseInstructor
		err         error
	}

	MockNatsPublisher struct {
		publishedEvents []string
	}
)

func (m *MockCourseRepository) Create(ctx context.Context, course *domain.Course) error {
	if m.err != nil {
		return m.err
	}
	if m.courses == nil {
		m.courses = make(map[uuid.UUID]*domain.Course)
	}
	m.courses[course.ID] = course
	return nil
}

func (m *MockCourseRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Course, error) {
	if m.err != nil {
		return nil, m.err
	}
	course, exists := m.courses[id]
	if !exists {
		return nil, fmt.Errorf("not found")
	}
	return course, nil
}

func (m *MockCourseRepository) Update(ctx context.Context, id uuid.UUID, updateFn func(item *domain.Course) (bool, error)) error {
	course, exists := m.courses[id]
	if !exists {
		return fmt.Errorf("not found")
	}
	updated, err := updateFn(course)
	if err != nil {
		return err
	}
	if updated {
		course.UpdatedAt = time.Now().UTC()
	}
	return nil
}

func (m *MockCourseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	if _, exists := m.courses[id]; !exists {
		return fmt.Errorf("not found")
	}
	delete(m.courses, id)
	return nil
}

func (m *MockCourseRepository) List(ctx context.Context) ([]*domain.Course, error) {
	if m.err != nil {
		return nil, m.err
	}
	courses := make([]*domain.Course, 0, len(m.courses))
	for _, c := range m.courses {
		courses = append(courses, c)
	}
	return courses, nil
}

func (m *MockGroupRepository) Create(ctx context.Context, course *domain.Group) error {
	return nil
}

func (m *MockGroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
	if m.err != nil {
		return nil, m.err
	}
	return nil, fmt.Errorf("error getByID")
}

func (m *MockGroupRepository) Update(ctx context.Context, id uuid.UUID, updateFn func(item *domain.Group) (bool, error)) error {
	return nil
}

func (m *MockGroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.err != nil {
		return m.err
	}

	return nil
}

func (m *MockGroupRepository) List(ctx context.Context) ([]*domain.Group, error) {
	if m.err != nil {
		return nil, m.err
	}
	courses := make([]*domain.Group, 0, len(m.groups))
	for _, c := range m.groups {
		courses = append(courses, c)
	}
	return courses, nil
}

func (m *MockGroupRepository) ListByCourse(context.Context, uuid.UUID) ([]*domain.Group, error) {
	return nil, fmt.Errorf("list by course error")
}

// Member repository mock methods
func (m *MockMemberRepository) Create(ctx context.Context, member *domain.GroupMember) error {
	if m.err != nil {
		return m.err
	}
	m.members = append(m.members, member)
	return nil
}

func (m *MockMemberRepository) Delete(ctx context.Context, groupID, userID uuid.UUID, role domain.MemberRole) error {
	if m.err != nil {
		return m.err
	}
	for i, member := range m.members {
		if member.GroupID == groupID && member.UserID == userID && member.Role == role {
			m.members = append(m.members[:i], m.members[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("not found")
}

func (m *MockMemberRepository) ListByGroup(ctx context.Context, groupID uuid.UUID) ([]*domain.GroupMember, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []*domain.GroupMember
	for _, member := range m.members {
		if member.GroupID == groupID {
			result = append(result, member)
		}
	}
	return result, nil
}

func (m *MockMemberRepository) Exists(ctx context.Context, groupID, userID uuid.UUID, role domain.MemberRole) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	for _, member := range m.members {
		if member.GroupID == groupID && member.UserID == userID && member.Role == role {
			return true, nil
		}
	}
	return false, nil
}

// Instructor repository mock methods (similar pattern)
// ... [Implement similar methods for InstructorRepository] ...

// Nats publisher mock methods
func (m *MockNatsPublisher) PublishCourseCreated(ctx context.Context, course *domain.Course) {
	m.publishedEvents = append(m.publishedEvents, "course_created")
}

func (m *MockNatsPublisher) PublishGroupCreated(ctx context.Context, group *domain.Group) {
	m.publishedEvents = append(m.publishedEvents, "group_created")
}

func (m *MockNatsPublisher) PublishGroupMemberAdded(ctx context.Context, member *domain.GroupMember) {
	m.publishedEvents = append(m.publishedEvents, "member_added")
}

func (m *MockNatsPublisher) PublishCourseInstructorAssigned(ctx context.Context, instructor *domain.CourseInstructor) {
	m.publishedEvents = append(m.publishedEvents, "instructor_assigned")
}

// Test helpers
func newTestService() *service.ManagersService {
	return service.NewManagersService(
		func() time.Time { return time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC) },
		nil, // logger
		&MockCourseRepository{
			courses: make(map[uuid.UUID]*domain.Course),
		},
		&MockGroupRepository{
			groups: make(map[uuid.UUID]*domain.Group),
		},
		&MockMemberRepository{},
		nil,
		nil,
		nil,
	)
}

// Tests
func TestCreateCourse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := newTestService()
		course, err := s.CreateCourse(context.Background(), "Math 101")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if course.Name != "Math 101" {
			t.Errorf("expected name 'Math 101', got '%s'", course.Name)
		}
		if len(course.ID) == 0 {
			t.Error("expected non-empty UUID")
		}
	})
}

func TestUpdateCourse(t *testing.T) {
	s := newTestService()
	course, _ := s.CreateCourse(context.Background(), "Initial Name")

	t.Run("success", func(t *testing.T) {
		newName := "Updated Name"
		updated, err := s.UpdateCourse(context.Background(), domain.CourseUpdate{
			ID:   course.ID,
			Name: &newName,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.Name != newName {
			t.Errorf("expected name '%s', got '%s'", newName, updated.Name)
		}
	})

	t.Run("invalid name", func(t *testing.T) {
		invalidName := "A"
		_, err := s.UpdateCourse(context.Background(), domain.CourseUpdate{
			ID:   course.ID,
			Name: &invalidName,
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestCreateGroup(t *testing.T) {
	s := newTestService()
	course, _ := s.CreateCourse(context.Background(), "Physics 101")

	t.Run("success", func(t *testing.T) {
		group, err := s.CreateGroup(context.Background(), course.ID, "Group A")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if group.Name != "Group A" {
			t.Errorf("expected name 'Group A', got '%s'", group.Name)
		}
		if group.CourseID != course.ID {
			t.Errorf("expected course ID %s, got %s", course.ID, group.CourseID)
		}
	})

	t.Run("invalid name", func(t *testing.T) {
		_, err := s.CreateGroup(context.Background(), course.ID, "A")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("course not found", func(t *testing.T) {
		_, err := s.CreateGroup(context.Background(), uuid.New(), "Valid Name")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestAddGroupMember(t *testing.T) {
	s := newTestService()
	userID := uuid.New()
	t.Run("group not found", func(t *testing.T) {
		_, err := s.AddGroupMember(context.Background(), uuid.New(), userID, domain.TeacherRole)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestAssignCourseInstructor(t *testing.T) {
	s := newTestService()
	userID := uuid.New()
	t.Run("course not found", func(t *testing.T) {
		_, err := s.AssignCourseInstructor(context.Background(), uuid.New(), userID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
