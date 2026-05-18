package server

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "dsp-stream-engine/pkg/api/stream"
)

// setupServer sets up an in-memory gRPC server for testing.
func setupServer() (*grpc.Server, *bufconn.Listener) {
	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	pb.RegisterDSPStreamServer(s, NewDSPStreamServer())
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
	return s, lis
}

// dialer creates a client connection to the in-memory server.
func dialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
}

func TestDSPStreamServer_ProcessAudio(t *testing.T) {
	srv, lis := setupServer()
	defer srv.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer(lis)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewDSPStreamClient(conn)
	stream, err := client.ProcessAudio(ctx)
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}

	// 8 bytes -> 2 float32s: 1.0, 1.0
	// 1.0 in IEEE 754 float32 is 0x3f800000 -> 0x00, 0x00, 0x80, 0x3f
	inputData := []byte{0x00, 0x00, 0x80, 0x3f, 0x00, 0x00, 0x80, 0x3f}

	err = stream.Send(&pb.AudioChunk{Data: inputData})
	if err != nil {
		t.Fatalf("Failed to send chunk: %v", err)
	}

	err = stream.CloseSend()
	if err != nil {
		t.Fatalf("Failed to close send: %v", err)
	}

	res, err := stream.Recv()
	if err != nil && err != io.EOF {
		t.Fatalf("Failed to receive chunk: %v", err)
	}

	if len(res.GetData()) != len(inputData) {
		t.Errorf("Expected length %d, got %d", len(inputData), len(res.GetData()))
	}

	// Server applies a simple moving average filter: y[n] = 0.5*x[n] + 0.5*x[n-1]
	// First input 1.0 -> output 0.5 (0x3f000000 -> 0x00, 0x00, 0x00, 0x3f)
	expectedFirstOutput := []byte{0x00, 0x00, 0x00, 0x3f}
	for i := 0; i < 4; i++ {
		if res.GetData()[i] != expectedFirstOutput[i] {
			t.Errorf("Mismatch at byte %d: got %x, expected %x", i, res.GetData()[i], expectedFirstOutput[i])
		}
	}
}
