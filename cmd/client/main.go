package main

import (
	"context"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "dsp-stream-engine/pkg/api/stream"
)

// HeaderSize defines the size of the custom IEEE Float WAV header.
const HeaderSize = 56

// main reads an input WAV file, streams it via gRPC for DSP processing, and writes the output.
func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <input.wav> <output.wav>", os.Args[0])
	}
	inputFile := os.Args[1]
	outputFile := os.Args[2]

	in, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Failed to open input: %v", err)
	}
	defer in.Close()

	out, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create output: %v", err)
	}
	defer out.Close()

	header := make([]byte, HeaderSize)
	if _, err := io.ReadFull(in, header); err != nil {
		log.Fatalf("Failed to read header: %v", err)
	}
	if _, err := out.Write(header); err != nil {
		log.Fatalf("Failed to write header: %v", err)
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewDSPStreamClient(conn)

	stream, err := client.ProcessAudio(context.Background())
	if err != nil {
		log.Fatalf("could not open stream: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive processed chunk: %v", err)
			}
			if _, err := out.Write(res.GetData()); err != nil {
				log.Fatalf("Failed to write processed data: %v", err)
			}
		}
	}()

	buf := make([]byte, 4096)
	for {
		n, err := in.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to read chunk: %v", err)
		}
		if err := stream.Send(&pb.AudioChunk{Data: buf[:n]}); err != nil {
			log.Fatalf("Failed to send chunk: %v", err)
		}
	}
	stream.CloseSend()
	<-waitc
	log.Println("Processing complete. Saved to", outputFile)
}
