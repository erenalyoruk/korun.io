FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ./go.work .

COPY ./auth-service/go.mod ./auth-service/go.mod
COPY ./auth-service/go.sum ./auth-service/go.sum

COPY ./shared/go.mod ./shared/go.mod
COPY ./shared/go.sum ./shared/go.sum

RUN go mod download

COPY ./auth-service ./auth-service
COPY ./shared ./shared

RUN CGO_ENABLED=0 GOOS=linux go build -o ./service ./auth-service/cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root

COPY --from=builder /app/service .
COPY ./auth-service/configs ./configs

EXPOSE 8081

CMD ["./service"]
