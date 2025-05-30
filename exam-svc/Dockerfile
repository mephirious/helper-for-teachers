FROM golang:1.24.0-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o exam-svc ./cmd/exam/main.go

CMD ["./exam-svc"]
