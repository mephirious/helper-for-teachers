#!/bin/bash

GRPC_ADDR=${GRPC_ADDR:-"localhost:50051"}

# Check grpcurl
if ! command -v grpcurl &> /dev/null; then
    echo "❌ grpcurl not found. Install it with: go install github.com/fullstorydev/grpcurl/...@latest"
    exit 1
fi

echo "🚀 Testing ExamService endpoints at $GRPC_ADDR..."

# Function to call endpoint
call_endpoint() {
    local method=$1
    local payload=$2
    echo "➡️ Calling $method..."
    if [ -z "$payload" ]; then
        grpcurl -plaintext "$GRPC_ADDR" "examservice.ExamService/$method"
    else
        echo "$payload" | grpcurl -plaintext -d @ "$GRPC_ADDR" "examservice.ExamService/$method"
    fi
    echo ""
}

TIMESTAMP=$(date +%s)

# Payloads
create_exam_payload=$(cat <<EOF
{
  "exam": {
    "title": "Math Exam",
    "description": "A test on basic algebra.",
    "created_by": "507f191e810c19729de860ea",
    "status": "draft",
    "created_at": { "seconds": $TIMESTAMP },
    "updated_at": { "seconds": $TIMESTAMP }
  }
}
EOF
)

create_task_payload=$(cat <<EOF
{
  "task": {
    "exam_id": "507f1f77bcf86cd799439011",
    "task_type": "writing",
    "description": "Write an essay on climate change.",
    "score": 10.0,
    "created_at": { "seconds": $TIMESTAMP }
  }
}
EOF
)

create_question_payload=$(cat <<EOF
{
  "question": {
    "exam_id": "507f1f77bcf86cd799439011",
    "question_text": "What is the capital of France?",
    "options": ["Paris", "London", "Berlin", "Madrid"],
    "correct_answer": "Paris",
    "status": "active",
    "created_at": { "seconds": $TIMESTAMP }
  }
}
EOF
)

generate_exam_payload=$(cat <<EOF
{
  "user_id": "507f191e810c19729de860ea",
  "num_questions": 5,
  "num_tasks": 3,
  "topic": "Algebra",
  "grade": "10th"
}
EOF
)

# Replace with real IDs after creating entities
TASK_ID="REPLACE_TASK_ID"
QUESTION_ID="REPLACE_QUESTION_ID"
EXAM_ID="REPLACE_EXAM_ID"

# Calls
call_endpoint "CreateExam" "$create_exam_payload"
call_endpoint "GetAllExams"
call_endpoint "GetExamsByUser" '{"user_id": "507f191e810c19729de860ea"}'
call_endpoint "GetExamByID" "{\"id\": \"$EXAM_ID\"}"
call_endpoint "UpdateExam" "$create_exam_payload"
call_endpoint "UpdateExamStatus" "{\"id\": \"$EXAM_ID\", \"status\": \"archived\"}"
call_endpoint "DeleteExam" "{\"id\": \"$EXAM_ID\"}"
call_endpoint "GetExamWithDetails" "{\"id\": \"$EXAM_ID\"}"
call_endpoint "GenerateExamUsingAI" "$generate_exam_payload"

call_endpoint "CreateTask" "$create_task_payload"
call_endpoint "GetAllTasks"
call_endpoint "GetTaskByID" "{\"id\": \"$TASK_ID\"}"
call_endpoint "GetTasksByExamID" "{\"exam_id\": \"507f1f77bcf86cd799439011\"}"
call_endpoint "UpdateTask" "$create_task_payload"
call_endpoint "DeleteTask" "{\"id\": \"$TASK_ID\"}"

call_endpoint "CreateQuestion" "$create_question_payload"
call_endpoint "GetAllQuestions"
call_endpoint "GetQuestionByID" "{\"id\": \"$QUESTION_ID\"}"
call_endpoint "GetQuestionsByExamID" "{\"exam_id\": \"507f1f77bcf86cd799439011\"}"
call_endpoint "UpdateQuestion" "$create_question_payload"
call_endpoint "DeleteQuestion" "{\"id\": \"$QUESTION_ID\"}"

echo "✅ All endpoints called."
