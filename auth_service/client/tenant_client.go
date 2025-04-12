package client

import (
	"log"

	"github.com/krishnachowdaryvanam/authboard/tenant_service/tenantpb"
	"google.golang.org/grpc"
)

var tenantClient tenantpb.TenantServiceClient

// InitTenantClient initializes the gRPC connection to the tenant service
func InitTenantClient() {
	conn, err := grpc.Dial("tenant-service:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to tenant service: %v", err)
	}
	tenantClient = tenantpb.NewTenantServiceClient(conn)
}

func GetTenantClient() tenantpb.TenantServiceClient {
	return tenantClient
}
