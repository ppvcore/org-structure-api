FROM golang:1.26-alpine AS builder
WORKDIR /test
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/app
FROM alpine:latest
WORKDIR /test
COPY --from=builder --chown=root:root --chmod=755 /test/main .
COPY --from=builder /test/.env .
EXPOSE 8080
CMD ["./main"]