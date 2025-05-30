GRPC_ADDR=${GRPC_ADDR:-"localhost:4001"}
TMP_DIR="/tmp/exam_ids"
mkdir -p "$TMP_DIR"

if ! command -v grpcurl &> /dev/null; then
  echo "❌ grpcurl не найден. Установи: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
  exit 1
fi

if ! command -v jq &> /dev/null; then
  echo "❌ jq не найден. Установи: sudo apt install jq"
  exit 1
fi

TIMESTAMP=$(date +%s)

call() {
  local method=$1
  local data=$2
  echo "➡️  $method"
  if [ -z "$data" ]; then
    grpcurl -plaintext "$GRPC_ADDR" "examservice.ExamService/$method"
  else
    echo "$data" | grpcurl -plaintext -d @ "$GRPC_ADDR" "examservice.ExamService/$method"
  fi
}

update_exam_payload=$(cat <<EOF
{
  "exam": {
    "id": "$EXAM_ID",
    "title": "Updated Math Exam",
    "description": "An updated test on algebra and geometry.",
    "created_by": "507f191e810c19729de860ea",
    "status": "active",
    "created_at": { "seconds": $TIMESTAMP },
    "updated_at": { "seconds": $TIMESTAMP }
  }
}
EOF
)

case "$1" in
  1) 
    payload=$(cat <<EOF
{
  "exam": {
    "title": "Math Exam",
    "description": "Algebra basics",
    "created_by": "507f191e810c19729de860ea",
    "status": "draft"
  }
}
EOF
)
    result=$(echo "$payload" | grpcurl -plaintext -d @ "$GRPC_ADDR" examservice.ExamService/CreateExam)
    echo "$result"
    echo "$result" | jq -r '.id' > "$TMP_DIR/exam_id"
    ;;

  2) call GetAllExams ;;
  3) call GetExamsByUser '{"user_id":"507f191e810c19729de860ea"}' ;;
  4) call GetExamByID "{\"id\":\"6838ecdf6b92918b4634882a\"}" ;;
  5) call UpdateExam "$(cat <<EOF
{
  "exam": {
    "id": "6838ecdf6b92918b4634882a",
    "title": "Updated Math Exam",
    "description": "An updated test on algebra and geometry.",
    "created_by": "507f191e810c19729de860ea",
    "status": "active"
  }
}
EOF
)" ;;

  6) call UpdateExamStatus "{\"id\":\"6838ecdf6b92918b4634882a\", \"status\":\"archived\"}" ;;
  7) call DeleteExam "{\"id\":\"6838ecdf6b92918b4634882a\"}" ;;
  8) call GetExamWithDetails "{\"id\":\"6838e7377204d91e28c005c8\"}" ;;
  9) call GenerateExamUsingAI "$(cat <<EOF
{
  "user_id": "507f191e810c19729de860ea",
  "num_questions": 3,
  "num_tasks": 2,
  "topic": "Geometry",
  "grade": "10th"
}
EOF
)" ;;

  10) 
    payload=$(cat <<EOF
{
  "task": {
    "exam_id": "6838f0f672177a8bb05d347d",
    "task_type": "writing",
    "description": "Write an essay.",
    "score": 10.0
  }
}
EOF
)
    result=$(echo "$payload" | grpcurl -plaintext -d @ "$GRPC_ADDR" examservice.ExamService/CreateTask)
    echo "$result"
    echo "$result" | jq -r '.id' > "$TMP_DIR/task_id"
    ;;

  11) call GetAllTasks ;;
  12) call GetTaskByID "{\"id\":\"6838f2b05293c5a90799ac99\"}" ;;
  13) call GetTasksByExamID "{\"exam_id\":\"6838f2b05293c5a90799ac95\"}" ;;
  14) call UpdateTask "$(cat <<EOF
{
  "task": {
    "id": "6838f2b05293c5a90799ac99",
    "exam_id": "6838f2b05293c5a90799ac95",
    "task_type": "short_answer",
    "description": "Updated task desc",
    "score": 15.0
  }
}
EOF
)" ;;
  15) call DeleteTask "{\"id\":\"6838f2b05293c5a90799ac9a\"}" ;;

  16)
    payload=$(cat <<EOF
{
  "question": {
    "exam_id": "6838f2b05293c5a90799ac95",
    "question_text": "Capital of France?",
    "options": ["Paris","London","Berlin","Madrid"],
    "correct_answer": "Paris",
    "status": "active"
  }
}
EOF
)
    result=$(echo "$payload" | grpcurl -plaintext -d @ "$GRPC_ADDR" examservice.ExamService/CreateQuestion)
    echo "$result"
    echo "$result" | jq -r '.id' > "$TMP_DIR/question_id"
    ;;

  17) call GetAllQuestions ;;
  18) call GetQuestionByID "{\"id\":\"6838f0f672177a8bb05d3484\"}" ;;
  19) call GetQuestionsByExamID "{\"exam_id\":\"6838f0f672177a8bb05d347d\"}" ;;
  20) call UpdateQuestion "$(cat <<EOF
{
  "question": {
    "id": "6838f0f672177a8bb05d3484",
    "exam_id": "6838f0f672177a8bb05d347d",
    "question_text": "Capital of Germany?",
    "options": ["Berlin","Paris","Rome","Lisbon"],
    "correct_answer": "Berlin",
    "status": "active"
  }
}
EOF
)" ;;
  21) call DeleteQuestion "{\"id\":\"6838f0f672177a8bb05d3484\"}" ;;

  *)
    echo "❓ Unknown or missing command. Use number 1-21."
    ;;
esac
