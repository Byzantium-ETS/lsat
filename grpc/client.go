package grpc

import (
	"log"
	"lsat/challenge"

	"google.golang.org/grpc"
	// Import any other necessary packages here
)

var opts []grpc.DialOption

func InitGrpcClient(address string) challenge.LightningNode {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client using the connection
	client := challenge.NewLndClient(conn)

	return client
}
