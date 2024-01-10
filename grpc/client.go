package grpc

import (
	"log"
	"lsat/challenge"

	"google.golang.org/grpc"
	// Import any other necessary packages here
)

func InitGrpcClient(address string) challenge.LightningNode {
	// Set up a connection to the gRPC server.
	// Replace "localhost:50051" with the actual server address and port.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client using the connection
	client := challenge.NewLndClient(conn)

	return client
}
