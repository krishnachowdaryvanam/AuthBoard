FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o user_service .

# Install netcat for the wait script
RUN apt-get update && apt-get install -y netcat-openbsd

# Expose gRPC port
EXPOSE 50052

# Use wait script to delay start until DB is ready
CMD ["./wait-for-postgres.sh", "./user_service"]