FROM golang:1.18 AS development
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/cespare/reflex@latest
EXPOSE 8080
CMD reflex -sr '\.go$' go run cmd/server/main.go
