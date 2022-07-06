FROM golang:1.18 AS development
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/cespare/reflex@latest
EXPOSE 8080
CMD reflex -sr '\.go$' go run cmd/server/main.go

FROM golang:1.18 AS builder
ENV GOOS linux
ENV CGO_ENABLED 0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app

FROM alpine:latest AS production
RUN apk add --no-cache ca-certificates
COPY --from=builder app .
EXPOSE 8080
CMD ./app
