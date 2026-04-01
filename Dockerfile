FROM golang:1.25.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o comment-system ./server.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/comment-system /app/comment-system

EXPOSE 8080

CMD ["/app/comment-system"]