package main

import (
	"fmt"
	"log"
	"net"
	"os"

	userpb "github.com/krishnachowdaryvanam/authboard/user_service/userspb"
	"google.golang.org/grpc"
)

func main() {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	if err := InitDB(connStr); err != nil {
		log.Fatalf("Failed to connect to user DB: %v", err)
	}

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &UserServer{})

	log.Println("User gRPC service running on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
