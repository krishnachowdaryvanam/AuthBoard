package client

import (
	"log"

	"github.com/krishnachowdaryvanam/authboard/proto/eventpb"
	"google.golang.org/grpc"
)

var eventClient eventpb.EventServiceClient

// InitEventClient initializes the gRPC connection to the event service
func InitEventClient() {
	conn, err := grpc.Dial("event-service:50054", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to event service: %v", err)
	}
	eventClient = eventpb.NewEventServiceClient(conn)
}

// GetEventClient returns the initialized EventServiceClient
func GetEventClient() eventpb.EventServiceClient {
	return eventClient
}
