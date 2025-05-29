package smtp

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
	"time"

	"github.com/mailersend/mailersend-go"
	eventspb "github.com/suyundykovv/margulan-protos/gen/go/events/v1"
	"go.uber.org/zap"
)

const (
	fromEmail    = "notifications@yourschool.edu"
	fromName     = "School Notification System"
	templatePath = "./templates/notification.html"
)

type MailerSendConfig struct {
	APIKey        string
	Domain        string
	SenderEmail   string
	SenderName    string
	TemplatePath  string
	WebhookSecret string
}

type EmailNotification struct {
	RecipientName  string
	RecipientEmail string
	ItemType       string 
	Title          string
	Description    string
	Time           string
	Location       string
	ActionURL      string
}

type Mailer struct {
	client   *mailersend.Mailersend
	logger   *zap.Logger
	config   MailerSendConfig
	emailTpl *template.Template
}

func NewMailer(config MailerSendConfig, logger *zap.Logger) (*Mailer, error) {
	ms := mailersend.NewMailersend(config.APIKey)

	tpl, err := template.ParseFiles(config.TemplatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse email template: %w", err)
	}

	return &Mailer{
		client:   ms,
		logger:   logger,
		config:   config,
		emailTpl: tpl,
	}, nil
}

func (m *Mailer) SendLessonNotification(ctx context.Context, lesson *eventspb.Lesson, recipient *eventspb.User) error {
	location := "Online"
	if !lesson.IsOnline && lesson.Classroom != "" {
		location = lesson.Classroom
	}

	notification := EmailNotification{
		RecipientName:  recipient.Name,
		RecipientEmail: recipient.Email,
		ItemType:       "lesson",
		Title:          lesson.Title,
		Description:    lesson.Description,
		Time:           lesson.StartTime.AsTime().Format("January 2, 2006 15:04"),
		Location:       location,
		ActionURL:      fmt.Sprintf("https://school.edu/lessons/%s", lesson.Id),
	}

	return m.sendNotification(ctx, notification)
}

// SendTaskNotification sends email notification for upcoming task
func (m *Mailer) SendTaskNotification(ctx context.Context, task *eventspb.Task, recipient *eventspb.User) error {
	notification := EmailNotification{
		RecipientName:  recipient.Name,
		RecipientEmail: recipient.Email,
		ItemType:       "task",
		Title:          task.Title,
		Description:    task.Description,
		Time:           task.DueDate.AsTime().Format("January 2, 2006 15:04"),
		Location:       fmt.Sprintf("Course: %s", task.CourseId),
		ActionURL:      fmt.Sprintf("https://school.edu/tasks/%s", task.Id),
	}

	return m.sendNotification(ctx, notification)
}

func (m *Mailer) sendNotification(ctx context.Context, data EmailNotification) error {
	// Prepare email body from template
	var body string
	buf := new(bytes.Buffer)
	if err := m.emailTpl.Execute(buf, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}
	body = buf.String()

	// Create Mailersend message
	message := m.client.Email.NewMessage()
	message.SetFrom(m.config.SenderEmail, m.config.SenderName)
	message.SetRecipients([]mailersend.Recipient{
		{
			Email: data.RecipientEmail,
			Name:  data.RecipientName,
		},
	})
	message.SetSubject(fmt.Sprintf("Upcoming %s: %s", data.ItemType, data.Title))
	message.SetHTML(body)
	message.SetText(fmt.Sprintf(
		"Reminder for %s: %s\nTime: %s\nLocation: %s\nDescription: %s\nView: %s",
		data.ItemType,
		data.Title,
		data.Time,
		data.Location,
		data.Description,
		data.ActionURL,
	))

	// Send email
	res, err := m.client.Email.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	m.logger.Info("Email notification sent",
		zap.String("recipient", data.RecipientEmail),
		zap.String("type", data.ItemType),
		zap.String("message_id", res.Header.Get("X-Message-Id")))

	return nil
}

// ScheduleNotifications sets up timed notifications for lessons and tasks
func (m *Mailer) ScheduleNotifications(ctx context.Context, scheduler SchedulerService) {
	go func() {
		ticker := time.NewTicker(15 * time.Minute) // Check every 15 minutes
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.checkUpcomingItems(ctx, scheduler)
			}
		}
	}()
}

func (m *Mailer) checkUpcomingItems(ctx context.Context, scheduler SchedulerService) {
	now := time.Now()
	windowStart := now.Add(25 * time.Minute) // 25 minutes from now
	windowEnd := now.Add(24 * time.Hour)     // Up to 24 hours ahead

	// Get upcoming lessons
	lessons, err := scheduler.GetUpcomingLessons(ctx, windowStart, windowEnd)
	if err != nil {
		m.logger.Error("Failed to get upcoming lessons", zap.Error(err))
		return
	}

	// Get upcoming tasks
	tasks, err := scheduler.GetUpcomingTasks(ctx, windowStart, windowEnd)
	if err != nil {
		m.logger.Error("Failed to get upcoming tasks", zap.Error(err))
		return
	}

	// Process notifications
	m.sendNotificationsForItems(ctx, lessons, tasks)
}

func (m *Mailer) sendNotificationsForItems(ctx context.Context, lessons []*eventspb.Lesson, tasks []*eventspb.Task) {
	// Track sent notifications to avoid duplicates
	sent := make(map[string]bool)

	// Process lessons
	for _, lesson := range lessons {
		key := fmt.Sprintf("lesson:%s", lesson.Id)
		if sent[key] {
			continue
		}

		// Get recipients (students in group)
		recipients, err := m.getLessonRecipients(ctx, lesson.GroupId)
		if err != nil {
			m.logger.Error("Failed to get lesson recipients",
				zap.String("lesson_id", lesson.Id),
				zap.Error(err))
			continue
		}

		// Send to each recipient
		for _, recipient := range recipients {
			if err := m.SendLessonNotification(ctx, lesson, recipient); err != nil {
				m.logger.Error("Failed to send lesson notification",
					zap.String("lesson_id", lesson.Id),
					zap.String("recipient", recipient.Email),
					zap.Error(err))
				continue
			}
			sent[key] = true
		}
	}

	// Process tasks
	for _, task := range tasks {
		key := fmt.Sprintf("task:%s", task.Id)
		if sent[key] {
			continue
		}

		// Get recipients (students in group)
		recipients, err := m.getTaskRecipients(ctx, task.GroupId)
		if err != nil {
			m.logger.Error("Failed to get task recipients",
				zap.String("task_id", task.Id),
				zap.Error(err))
			continue
		}

		// Send to each recipient
		for _, recipient := range recipients {
			if err := m.SendTaskNotification(ctx, task, recipient); err != nil {
				m.logger.Error("Failed to send task notification",
					zap.String("task_id", task.Id),
					zap.String("recipient", recipient.Email),
					zap.Error(err))
				continue
			}
			sent[key] = true
		}
	}
}

// Helper methods would interface with your user/group services
func (m *Mailer) getLessonRecipients(ctx context.Context, groupID string) ([]*eventspb.User, error) {
	// Implement actual logic to get students in a group
	return []*eventspb.User{}, nil
}

func (m *Mailer) getTaskRecipients(ctx context.Context, groupID string) ([]*eventspb.User, error) {
	// Implement actual logic to get students in a group
	return []*eventspb.User{}, nil
}
