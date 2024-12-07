FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o main ./cmd/main.go

EXPOSE 3002
CMD ["./main"]


# docker system prune -af
# docker volume prune -f
