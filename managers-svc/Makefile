run:
	bin/managers

build:
	go build -o bin/managers cmd/app/main.go

proto:
	protoc   -I=internal/grpc/  --go_out=paths=source_relative:internal/grpc   --go-grpc_out=paths=source_relative:internal/grpc   internal/grpc/managers.proto
	protoc   -I=internal/grpc/  --go_out=paths=source_relative:internal/grpc   --go-grpc_out=paths=source_relative:internal/grpc   internal/grpc/events.proto