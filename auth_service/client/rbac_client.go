package client

import (
	"log"

	"github.com/krishnachowdaryvanam/authboard/rbac_service/rbacpb"
	"google.golang.org/grpc"
)

var rbacClient rbacpb.RbacServiceClient

// InitRBACClient initializes the gRPC connection to the RBAC service
func InitRBACClient() {
	conn, err := grpc.Dial("rbac-service:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to RBAC service: %v", err)
	}
	rbacClient = rbacpb.NewRbacServiceClient(conn)
}

func GetRBACClient() rbacpb.RbacServiceClient {
	return rbacClient
}
