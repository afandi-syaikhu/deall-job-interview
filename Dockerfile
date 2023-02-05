# BUILD STAGE
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o interview-service main.go

# RUN STAGE
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/interview-service .
COPY --from=builder /app/config/dev-config.json ./config/config.json

EXPOSE 8080
CMD [ "/app/interview-service" ]
