package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/genai"
)

type Question struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	ExamID        primitive.ObjectID `json:"exam_id" bson:"exam_id"`
	QuestionText  string             `json:"question_text" bson:"question_text"`
	Options       []string           `json:"options" bson:"options"`
	CorrectAnswer string             `json:"correct_answer" bson:"correct_answer"`
	Status        string             `json:"status" bson:"status"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
}

type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	ExamID      primitive.ObjectID `json:"exam_id" bson:"exam_id"`
	TaskType    string             `json:"task_type" bson:"task_type"`
	Description string             `json:"description" bson:"description"`
	Score       float32            `json:"score" bson:"score"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type Exam struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	CreatedBy   primitive.ObjectID `json:"created_by" bson:"created_by"`
	Status      string             `json:"status" bson:"status"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type ExamGenResult struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Questions   []Question `json:"questions"`
	Tasks       []Task     `json:"tasks"`
}

func buildPrompt(numQuestions, numTasks int) string {
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
- Use realistic, clearly written English academic content (math).
- Ensure "correct_answer" matches one of the options in each question.
- Ensure "task_type" is one of: "writing", "matching", "short answer".
- Ensure "score" is a positive number.
`, numQuestions, numTasks)
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// Parse command-line arguments
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go <num_questions> <num_tasks>")
	}
	numQuestions, err := strconv.Atoi(os.Args[1])
	if err != nil || numQuestions < 0 {
		log.Fatalf("Invalid number of questions: %v\n", err)
	}
	numTasks, err := strconv.Atoi(os.Args[2])
	if err != nil || numTasks < 0 {
		log.Fatalf("Invalid number of tasks: %v\n", err)
	}

	// Initialize Gemini API client
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY not set in environment or .env file")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v\n", err)
	}
	//defer client.Close()

	// Build prompt and generate content
	prompt := buildPrompt(numQuestions, numTasks)

	config := &genai.GenerateContentConfig{
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

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-1.5-flash",
		genai.Text(prompt),
		config,
	)
	if err != nil {
		log.Fatalf("Failed to generate content: %v\n", err)
	}

	// Extract and clean JSON response
	rawJSON := result.Text()
	rawJSON = strings.TrimSpace(rawJSON)
	if strings.HasPrefix(rawJSON, "```json") {
		rawJSON = strings.TrimPrefix(rawJSON, "```json")
		rawJSON = strings.TrimSuffix(rawJSON, "```")
		rawJSON = strings.TrimSpace(rawJSON)
	}

	fmt.Println("🧪 Cleaned JSON:")
	fmt.Println(rawJSON)

	// Decode JSON into ExamGenResult
	var examResult ExamGenResult
	if err := json.Unmarshal([]byte(rawJSON), &examResult); err != nil {
		log.Fatalf("Error decoding JSON: %v\n", err)
	}

	// Create Exam with MongoDB ObjectIDs
	exam := Exam{
		ID:          primitive.NewObjectID(),
		Title:       examResult.Title,
		Description: examResult.Description,
		Status:      examResult.Status,
		CreatedBy:   primitive.NewObjectID(), // Placeholder for actual user
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Assign IDs and ExamID to Questions
	for i := range examResult.Questions {
		examResult.Questions[i].ID = primitive.NewObjectID()
		examResult.Questions[i].ExamID = exam.ID
		examResult.Questions[i].CreatedAt = time.Now()
	}

	// Assign IDs and ExamID to Tasks
	for i := range examResult.Tasks {
		examResult.Tasks[i].ID = primitive.NewObjectID()
		examResult.Tasks[i].ExamID = exam.ID
		examResult.Tasks[i].CreatedAt = time.Now()
	}

	// Output results
	fmt.Printf("✅ Exam:\n")
	fmt.Printf("Title       : %s\n", exam.Title)
	fmt.Printf("Description : %s\n", exam.Description)
	fmt.Printf("Status      : %s\n", exam.Status)
	fmt.Printf("Created By  : %s\n", exam.CreatedBy.Hex())
	fmt.Printf("Created At  : %s\n", exam.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()

	fmt.Printf("📝 Questions (%d):\n", len(examResult.Questions))
	for i, q := range examResult.Questions {
		fmt.Printf("  %d. %s\n", i+1, q.QuestionText)
		for j, opt := range q.Options {
			optionLabel := string('A' + j)
			fmt.Printf("     %s) %s\n", optionLabel, opt)
		}
		fmt.Printf("     ✅ Correct Answer: %s\n", q.CorrectAnswer)
		fmt.Println()
	}

	fmt.Printf("📋 Tasks (%d):\n", len(examResult.Tasks))
	for i, t := range examResult.Tasks {
		fmt.Printf("  %d. Type: %s\n", i+1, t.TaskType)
		fmt.Printf("     Description: %s\n", t.Description)
		fmt.Printf("     Score      : %.2f\n", t.Score)
		fmt.Println()
	}
}
