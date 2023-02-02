build:
	@echo " >> building interview-service binary"
	@go build -v -o interview-service main.go

run: build
	@./interview-service

migration-init:
	@go run migrations/main.go init

migration-up:
	@go run migrations/main.go up

migration-down:
	@go run migrations/main.go down

