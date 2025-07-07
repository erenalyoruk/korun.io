FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.work and all go.mod/go.sum files first
COPY go.work .
COPY auth-service/go.mod ./auth-service/
COPY auth-service/go.sum ./auth-service/
COPY shared/go.mod ./shared/
COPY shared/go.sum ./shared/

# Sync the workspace and download all dependencies
RUN go work sync

# Copy the rest of the source code
COPY auth-service ./auth-service
COPY shared ./shared

# Build the auth-service
RUN CGO_ENABLED=0 GOOS=linux go build -o ./service ./auth-service/cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root

COPY --from=builder /app/service .
COPY auth-service/configs ./configs

EXPOSE 8081

CMD ["./service"]