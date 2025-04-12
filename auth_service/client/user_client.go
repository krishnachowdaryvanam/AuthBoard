package client

import (
	"log"

	userpb "github.com/krishnachowdaryvanam/authboard/user_service/userspb"
	"google.golang.org/grpc"
)

var userClient userpb.UserServiceClient

// InitUserClient initializes the gRPC connection to the user service
func InitUserClient() {
	conn, err := grpc.Dial("user-service:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	userClient = userpb.NewUserServiceClient(conn)
}

// GetUserClient returns the initialized UserServiceClient
func GetUserClient() userpb.UserServiceClient {
	return userClient
}
