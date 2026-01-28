FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

CMD ["/main"]