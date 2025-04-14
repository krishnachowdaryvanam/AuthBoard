package main

import (
	"context"
	"encoding/json"
	"log"

	"net"

	"github.com/krishnachowdaryvanam/authboard/proto/eventpb"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// EventServiceServer implements the EventService gRPC service
type EventServiceServer struct {
	eventpb.UnimplementedEventServiceServer
}

// PublishEvent will publish an event to Kafka
func (s *EventServiceServer) PublishEvent(ctx context.Context, req *eventpb.EventRequest) (*eventpb.EventResponse, error) {
	// Create a Kafka writer to publish messages to a topic
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"}, // Update with your Kafka broker addresses
		Topic:    "event_topic",              // Kafka topic name
		Balancer: &kafka.LeastBytes{},
	})

	// Define the event payload as a map
	var eventData map[string]interface{}

	// Try to unmarshal the payload into a map to work with it easily
	err := json.Unmarshal([]byte(req.Payload), &eventData)
	if err != nil {
		log.Printf("Error unmarshaling event payload: %v", err)
		return &eventpb.EventResponse{Success: false}, err
	}

	// Logic to handle different event types
	switch req.EventType {
	case "user_created":
		// Handle user created event
		log.Printf("Handling user_created event: %v", eventData)
		// Add any custom logic if needed (e.g., enriching data)
	case "tenant_deleted":
		// Handle tenant deleted event
		log.Printf("Handling tenant_deleted event: %v", eventData)
		// Add any custom logic if needed
	default:
		log.Printf("Unknown event type: %s", req.EventType)
		return &eventpb.EventResponse{Success: false}, nil
	}

	// Create the Kafka message with the event payload
	message := kafka.Message{
		Value: []byte(req.Payload), // Send the original event payload as message value
	}

	// Publish the event to Kafka
	err = writer.WriteMessages(ctx, message)
	if err != nil {
		log.Printf("Error publishing event: %v", err)
		return &eventpb.EventResponse{Success: false}, err
	}

	log.Printf("Event %s published successfully", req.EventType)

	// Respond back to the client with success
	return &eventpb.EventResponse{Success: true}, nil
}

func main() {
	// Create a listener for gRPC server on port 50052
	listener, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Failed to listen on port 50052: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the EventServiceServer with the gRPC server
	eventpb.RegisterEventServiceServer(grpcServer, &EventServiceServer{})

	// Register reflection service on gRPC server (optional, for debugging)
	reflection.Register(grpcServer)

	// Start the gRPC server
	log.Println("Event Service is running on port 50052")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
