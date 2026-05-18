package server

import (
	"io"
	"log"
	"sync"
	"unsafe"

	"dsp-stream-engine/internal/dsp"
	pb "dsp-stream-engine/pkg/api/stream"
)

// DSPStreamServer handles gRPC streaming of audio data.
type DSPStreamServer struct {
	pb.UnimplementedDSPStreamServer
	bufferPool sync.Pool
}

// NewDSPStreamServer creates a new instance of DSPStreamServer.
func NewDSPStreamServer() *DSPStreamServer {
	return &DSPStreamServer{
		bufferPool: sync.Pool{
			New: func() interface{} {
				b := make([]byte, 4096)
				return &b
			},
		},
	}
}

// ProcessAudio receives audio chunks, processes them via C++ DSP, and sends them back.
func (s *DSPStreamServer) ProcessAudio(stream pb.DSPStream_ProcessAudioServer) error {
	filter := dsp.NewBiquadFilter(0.5, 0.5, 0, 0, 0)
	defer filter.Close()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error receiving from stream: %v", err)
			return err
		}

		data := req.GetData()
		numBytes := len(data)
		if numBytes == 0 {
			continue
		}

		numFloats := numBytes / 4
		if numFloats > 0 {
			floatSlice := unsafe.Slice((*float32)(unsafe.Pointer(&data[0])), numFloats)
			filter.Process(floatSlice)
		}

		res := &pb.AudioChunk{
			Data: data,
		}

		if err := stream.Send(res); err != nil {
			log.Printf("Error sending to stream: %v", err)
			return err
		}
	}
}
