package mailjet

import "fmt"

type TemplateType string

const (
	TaskVerifiedTemplate     TemplateType = "task_verified"
	QuestionVerifiedTemplate TemplateType = "question_verified"
	ExamVerifiedTemplate     TemplateType = "exam_verified"
)

type EmailTemplate struct {
	Subject  string
	TextBody string
	HTMLBody string
}

var templates = map[TemplateType]EmailTemplate{
	TaskVerifiedTemplate: {
		Subject:  "Task Verified",
		TextBody: "A task has been verified.",
		HTMLBody: "<h3>Your task has been verified.</h3>",
	},
	QuestionVerifiedTemplate: {
		Subject:  "Question Verified",
		TextBody: "A question has been verified.",
		HTMLBody: "<h3>Your question has been verified.</h3>",
	},
	ExamVerifiedTemplate: {
		Subject:  "Exam Verified",
		TextBody: "An exam has been verified.",
		HTMLBody: "<h3>Your exam has been verified.</h3>",
	},
}

func (m *MailjetClient) SendTemplateEmail(toEmail, toName string, templateType TemplateType) error {
	tmpl, ok := templates[templateType]
	if !ok {
		return fmt.Errorf("template type %s not found", templateType)
	}
	return m.SendEmail(toEmail, toName, tmpl.Subject, tmpl.TextBody, tmpl.HTMLBody)
}
