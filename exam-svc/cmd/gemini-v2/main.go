// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/adapter/gemini"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/domain"
	"google.golang.org/genai"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	ctx := context.Background()
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v\n", err)
	}

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

	genaiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v\n", err)
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY not set in environment or .env file")
	}

	client, err := gemini.NewClient(genaiClient, "gemini-1.5-flash")

	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v\n", err)
	}

	result, err := client.GenerateExam(ctx, numQuestions, numTasks, "5 grade", "basic math")
	if err != nil {
		log.Fatalf("Failed to generate exam: %v\n", err)
	}

	exam := createExam(result)
	questions := createQuestions(result.Questions, exam.ID)
	tasks := createTasks(result.Tasks, exam.ID)

	printExam(exam)
	printQuestions(questions)
	printTasks(tasks)
}

func createExam(result *gemini.ExamGenResult) domain.Exam {
	return domain.Exam{
		ID:          primitive.NewObjectID(),
		Title:       result.Title,
		Description: result.Description,
		Status:      result.Status,
		CreatedBy:   primitive.NewObjectID(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func createQuestions(geminiQuestions []domain.Question, examID primitive.ObjectID) []domain.Question {
	questions := make([]domain.Question, len(geminiQuestions))
	for i, q := range geminiQuestions {
		questions[i] = domain.Question{
			ID:            primitive.NewObjectID(),
			ExamID:        examID,
			QuestionText:  q.QuestionText,
			Options:       q.Options,
			CorrectAnswer: q.CorrectAnswer,
			Status:        q.Status,
			CreatedAt:     time.Now(),
		}
	}
	return questions
}

func createTasks(geminiTasks []domain.Task, examID primitive.ObjectID) []domain.Task {
	tasks := make([]domain.Task, len(geminiTasks))
	for i, t := range geminiTasks {
		tasks[i] = domain.Task{
			ID:          primitive.NewObjectID(),
			ExamID:      examID,
			TaskType:    t.TaskType,
			Description: t.Description,
			Score:       t.Score,
			CreatedAt:   time.Now(),
		}
	}
	return tasks
}

func printExam(exam domain.Exam) {
	fmt.Printf("Exam:\n")
	fmt.Printf("Title       : %s\n", exam.Title)
	fmt.Printf("Description : %s\n", exam.Description)
	fmt.Printf("Status      : %s\n", exam.Status)
	fmt.Printf("Created By  : %s\n", exam.CreatedBy.Hex())
	fmt.Printf("Created At  : %s\n\n", exam.CreatedAt.Format("2006-01-02 15:04:05"))
}

func printQuestions(questions []domain.Question) {
	fmt.Printf("Questions (%d):\n", len(questions))
	for i, q := range questions {
		fmt.Printf("  %d. %s\n", i+1, q.QuestionText)
		for j, opt := range q.Options {
			optionLabel := fmt.Sprint('A' + j)
			fmt.Printf("     %s) %s\n", optionLabel, opt)
		}
		fmt.Printf("     Correct Answer: %s\n\n", q.CorrectAnswer)
	}
}

func printTasks(tasks []domain.Task) {
	fmt.Printf("Tasks (%d):\n", len(tasks))
	for i, t := range tasks {
		fmt.Printf("  %d. Type: %s\n", i+1, t.TaskType)
		fmt.Printf("     Description: %s\n", t.Description)
		fmt.Printf("     Score      : %.2f\n\n", t.Score)
	}
}
