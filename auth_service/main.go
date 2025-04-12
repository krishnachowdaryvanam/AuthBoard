package main

import (
	"log"

	"github.com/krishnachowdaryvanam/authboard/auth_service/client"
	"github.com/krishnachowdaryvanam/authboard/auth_service/routers"
)

func main() {
	// Initialize gRPC clients
	client.InitUserClient()
	client.InitTenantClient()
	client.InitRBACClient()

	// Initialize the router
	r := routers.SetupRouter()

	// Start the HTTP server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start the server: ", err)
	}
}
