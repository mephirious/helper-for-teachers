package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"google.golang.org/genai"
)

type Client struct {
	client    *genai.Client
	modelName string
}

type ExamGenResult struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      string            `json:"status"`
	Questions   []domain.Question `json:"questions"`
	Tasks       []domain.Task     `json:"tasks"`
}

func NewClient(client *genai.Client, modelName string) (*Client, error) {
	return &Client{
		client:    client,
		modelName: modelName,
	}, nil
}

func (c *Client) GenerateExam(ctx context.Context, numQuestions, numTasks int, grade, topic string) (*ExamGenResult, error) {
	prompt := buildPrompt(numQuestions, numTasks, grade, topic)
	config := createGenerationConfig()

	result, err := c.client.Models.GenerateContent(
		ctx,
		c.modelName,
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	rawJSON := cleanJSONResponse(result.Text())

	var examResult ExamGenResult
	if err := json.Unmarshal([]byte(rawJSON), &examResult); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	return &examResult, nil
}

func buildPrompt(numQuestions, numTasks int, grade, topic string) string {
	return fmt.Sprintf(`
Generate a JSON object for an academic exam with the following structure:
{
  "title": "string",
  "description": "string",
  "status": "draft",
  "questions": [
    {
      "question_text": "string",
      "options": ["string", "string", "string", "string"],
      "correct_answer": "string",
      "status": "active"
    }
  ],
  "tasks": [
    {
      "task_type": "string",
      "description": "string",
      "score": number
    }
  ]
}

- Generate exactly %d questions and %d tasks.
- Use realistic, clearly written English academic content based on the topic: "%s".
- The content should be appropriate for grade: "%s".
- Ensure "correct_answer" matches one of the options in each question.
- Ensure "task_type" is one of: "writing", "matching", "short answer".
- Ensure "score" is a positive number.
`, numQuestions, numTasks, topic, grade)
}

func createGenerationConfig() *genai.GenerateContentConfig {
	return &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"title":       {Type: genai.TypeString},
				"description": {Type: genai.TypeString},
				"status":      {Type: genai.TypeString},
				"questions": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"question_text": {Type: genai.TypeString},
							"options": {
								Type:  genai.TypeArray,
								Items: &genai.Schema{Type: genai.TypeString},
							},
							"correct_answer": {Type: genai.TypeString},
							"status":         {Type: genai.TypeString},
						},
						PropertyOrdering: []string{"question_text", "options", "correct_answer", "status"},
					},
				},
				"tasks": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"task_type":   {Type: genai.TypeString},
							"description": {Type: genai.TypeString},
							"score":       {Type: genai.TypeNumber},
						},
						PropertyOrdering: []string{"task_type", "description", "score"},
					},
				},
			},
			PropertyOrdering: []string{"title", "description", "status", "questions", "tasks"},
		},
	}
}

func cleanJSONResponse(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimSuffix(raw, "```")
	return strings.TrimSpace(raw)
}
