FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY logger ./logger
COPY docs ./docs
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
