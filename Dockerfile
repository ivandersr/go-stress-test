FROM golang:1.24-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . /app
RUN go build -o main ./cmd/main.go
ENTRYPOINT ["/app/main", "load"]