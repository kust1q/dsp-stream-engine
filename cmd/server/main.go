package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"dsp-stream-engine/internal/server"
	pb "dsp-stream-engine/pkg/api/stream"
)

// main starts the gRPC server listening on port 50051.
func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	dspServer := server.NewDSPStreamServer()

	pb.RegisterDSPStreamServer(grpcServer, dspServer)

	log.Printf("Server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
