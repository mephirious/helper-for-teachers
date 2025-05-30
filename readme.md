# Helper for Teachers

`Helper for Teachers` is a service designed to assist educators in creating, managing, and verifying exams, questions, and tasks. The `exam-svc` service, written in Go, provides a gRPC-based API for managing educational content, integrated with MongoDB for data persistence, Redis for caching, NATS for event publishing, Gemini for AI-generated exams, and Mailjet for email notifications.

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Setup](#setup)
  - [Clone the Repository](#clone-the-repository)
  - [Environment Variables](#environment-variables)
  - [Docker Compose](#docker-compose)
  - [Build and Run exam-svc](#build-and-run-exam-svc)
- [Usage](#usage)
  - [gRPC Endpoints](#grpc-endpoints)
  - [Testing with test_endpoints.sh](#testing-with-test_endpointssh)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Overview
The `exam-svc` service is a microservice within the `Helper for Teachers` ecosystem. It allows teachers to:
- Create and manage exams, questions, and tasks.
- Update the status of exams, questions, and tasks, with email notifications sent via Mailjet when items are marked as "verified."
- Generate exams using AI through the Gemini API.
- Publish exam events to NATS for integration with other services.
- Cache frequently accessed data in Redis for performance.

The service uses MongoDB with a single-node replica set for transactional consistency, ensuring reliable data operations. It is deployed using Docker, with a `docker-compose.yml` configuration for MongoDB and initialization scripts.

## Features
- **Exam Management**: Create, update, delete, and retrieve exams with associated questions and tasks.
- **Question and Task Management**: Add, update, and delete questions and tasks linked to exams.
- **AI-Powered Exam Generation**: Generate exams with questions and tasks using the Gemini AI model.
- **Transactional Operations**: Ensure data consistency using MongoDB transactions for create and delete operations.
- **Email Notifications**: Send verification emails via Mailjet when exams, questions, or tasks are marked as "verified."
- **Event Publishing**: Publish exam-related events (create, update, delete) to NATS for event-driven architectures.
- **Caching**: Use Redis to cache exam, question, and task data for faster access.
- **gRPC API**: Expose a robust API for client applications to interact with the service.

## Prerequisites
- **Go**: Version 1.22 or later.
- **Docker**: For running MongoDB and other services.
- **Docker Compose**: For orchestrating containers.
- **MongoDB**: Configured as a replica set (handled via `docker-compose.yml`).
- **Environment Variables**: API keys and configuration for MongoDB, Redis, NATS, Gemini, and Mailjet.
- **gRPCurl**: For testing gRPC endpoints (optional, used in `test_endpoints.sh`).

## Setup

### Clone the Repository
```bash
git clone https://github.com/mephirious/helper-for-teachers.git
cd helper-for-teachers
```

### Environment Variables
Create a `.env` file in the root directory or set environment variables manually. Example:
```bash
MONGO_URI=mongodb://localhost:27027/exam_db?replicaSet=rs0
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
NATS_URL=nats://localhost:4222
GEMINI_API_KEY=your-gemini-api-key
GEMINI_MODEL=your-model-name
MAILJET_API_KEY=your-mailjet-api-key
MAILJET_SECRET_KEY=your-mailjet-secret-key
MAILJET_FROM_EMAIL=sender@example.com
MAILJET_FROM_NAME=Exam Service
MAILJET_ADMIN_EMAIL=admin@example.com
MAILJET_ADMIN_NAME=Admin
```

### Docker Compose
The `docker-compose.yml` file sets up a MongoDB instance with a single-node replica set. Ensure the following files are in place:
- `docker-compose.yml`: Configures the MongoDB container.
- `mongo-init.js`: Initializes the replica set (`rs0`).

Start the MongoDB service:
```bash
docker-compose up -d
```

Verify MongoDB is running and the replica set is initialized:
```bash
docker logs final-mongo
docker exec -it final-mongo mongosh
rs.status()
```
Look for `"stateStr": "PRIMARY"` in the `rs.status()` output.

### Build and Run exam-svc
Navigate to the `exam-svc` directory, build, and run the service:
```bash
cd exam-svc
go build -o exam-svc
./exam-svc
```
The service will start on port `50051` (configurable in `config.go`).

Alternatively, if running `exam-svc` in Docker, add it to `docker-compose.yml`:
```yaml
services:
  exam-svc:
    build: ./exam-svc
    ports:
      - "50051:50051"
    environment:
      - MONGO_URI=mongodb://final-mongo:27017/exam_db?replicaSet=rs0
      # Add other env vars
    networks:
      - mongo-net
```

## Usage

### gRPC Endpoints
The `exam-svc` service exposes gRPC endpoints defined in `proto/exam_service.proto`. Key endpoints include:
- `CreateExam`: Create a new exam.
- `GetExamByID`: Retrieve an exam by ID.
- `UpdateExamStatus`: Update an exam’s status (triggers email if set to "verified").
- `DeleteExam`: Delete an exam and its associated tasks/questions.
- `CreateQuestion` / `UpdateQuestion` / `DeleteQuestion`: Manage questions.
- `CreateTask` / `UpdateTask` / `DeleteTask`: Manage tasks.
- `GenerateExamUsingAI`: Generate an exam with AI.

See `proto/exam_service.proto` for the full API specification.

### Testing with test_endpoints.sh
A `test_endpoints.sh` script is provided to test gRPC endpoints using `grpcurl`. Example usage:
```bash
./test_endpoints.sh
```
The script tests endpoints like `CreateExam`, `CreateTask`, `CreateQuestion`, and `DeleteExam`. Ensure the exam ID exists before deletion (see [Troubleshooting](#troubleshooting)).

Example flow in `test_endpoints.sh`:
```bash
# Create an exam
create_exam_payload='{
    "exam": {
        "title": "Test Exam",
        "description": "This is a test exam",
        "created_by": "507f191e810c19729de860ea",
        "status": "draft",
        "created_at": { "seconds": '$(date +%s)' },
        "updated_at": { "seconds": '$(date +%s)' }
    }
}'
call_endpoint "CreateExam" "$create_exam_payload"
exam_id=$(cat /tmp/grpcurl_output.txt | grep '"id"' | awk -F'"' '{print $4}' | head -1)

# Delete the exam
delete_exam_payload="{\"id\": \"$exam_id\"}"
call_endpoint "DeleteExam" "$delete_exam_payload"
```

## Troubleshooting
- **MongoDB Replica Set Not Initialized**:
  - Error: `NamespaceNotFound` for `local.oplog.rs`.
  - Fix: Check `docker logs final-mongo` for `mongo-init.js` execution. Manually initialize:
    ```bash
    docker exec -it final-mongo mongosh
    rs.initiate({ _id: "rs0", members: [{ _id: 0, host: "final-mongo:27017" }] })
    ```
  - Clear volume if needed:
    ```bash
    docker-compose down
    docker volume rm helper-for-teachers_mongo-data
    docker-compose up -d
    ```

- **Session Error in DeleteExam**:
  - Error: `session was not created by this client`.
  - Fix: Ensure `MONGO_URI` includes `replicaSet=rs0` (e.g., `mongodb://localhost:27027/exam_db?replicaSet=rs0`). Verify single `mongo.Client` in `app.go`.

- **Exam Not Found in DeleteExam**:
  - Error: `DeleteExam` fails for a specific ID.
  - Fix: Ensure the exam exists:
    ```bash
    docker exec -it final-mongo mongosh
    use exam_db
    db.exams.find({ _id: ObjectId("6838e7377204d91e28c005c8") })
    ```
    Create the exam first using `CreateExam`.

- **gRPC Connection Issues**:
  - Ensure `exam-svc` is running on port `50051` and `grpcurl` is installed:
    ```bash
    grpcurl -plaintext localhost:50051 list
    ```

- **Email Notifications Not Sent**:
  - Verify Mailjet credentials in `.env`. Check logs for errors:
    ```bash
    tail -f exam-svc.log
    ```

## Contributing
Contributions are welcome! To contribute:
1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/your-feature`).
3. Commit changes (`git commit -m "Add your feature"`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Open a pull request.

Please include tests and update documentation as needed.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.