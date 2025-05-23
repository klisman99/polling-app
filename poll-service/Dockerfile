FROM golang:1.24

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /app/poll-service ./cmd/poll-service

EXPOSE 3002

CMD ["/app/poll-service"]