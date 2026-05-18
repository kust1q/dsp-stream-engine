# DSP Stream Engine
A high-performance, zero-copy real-time audio processing microservice combining Go's network concurrency with C++'s mathematical determinism.

# Technologies
[![Go](https://img.shields.io/badge/-Go-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![C++](https://img.shields.io/badge/-C++-00599C?style=flat-square&logo=c%2B%2B&logoColor=white)](https://isocpp.org/)
[![gRPC](https://img.shields.io/badge/-gRPC-244C5A?style=flat-square&logo=grpc&logoColor=white)](https://grpc.io/)
[![CMake](https://img.shields.io/badge/-CMake-064F8C?style=flat-square&logo=cmake&logoColor=white)](https://cmake.org/)
[![Protobuf](https://img.shields.io/badge/-Protobuf-000000?style=flat-square&logo=c&logoColor=white)](https://protobuf.dev/)

## Tech Stack

### Languages & Compilers
- **Go 1.22+**: Handles high-level networking, goroutine management, and memory pooling.
- **C++17**: Implements the mathematical DSP (Digital Signal Processing) core for deterministic, low-latency execution.
- **CGO**: Acts as the bridge between Go and C++.

### Networking & RPC
- **gRPC**: Bidirectional streaming protocol to transmit audio chunks with minimal overhead.
- **Protocol Buffers**: Binary serialization format for API contracts.

### DSP & Memory Management
- **Biquad IIR Filter**: A stateful 2nd-order infinite impulse response filter implemented in C++.
- **Zero-Copy Architecture**: Uses `unsafe.Pointer` to reinterpret incoming gRPC byte slices into `float32` arrays, passing them to C++ for *in-place* mutation, completely bypassing the Go Garbage Collector (GC) for audio data.
- **CMake**: Build system for generating the static C++ DSP library (`libdsp.a`).

## Project Structure

```
dsp-stream-engine/
├── cmd/                           # Executable binaries
│   ├── analyze/                   # Utility to read and compare WAV samples
│   ├── client/                    # gRPC client to stream .wav files
│   ├── genwav/                    # Utility to generate test .wav files
│   └── server/                    # gRPC server entry point
│
├── cpp_dsp/                       # C++ Mathematical Core
│   ├── CMakeLists.txt             # CMake build configuration
│   ├── dsp.cpp                    # Biquad filter implementation
│   └── dsp.h                      # C-API wrapper for CGO
│
├── internal/                      # Private application code
│   ├── dsp/                       # Go CGO bindings for the C++ library
│   └── server/                    # gRPC server logic and zero-copy handler
│
├── pkg/                           # Publicly importable packages
│   └── api/stream/                # Generated Protobuf Go code
│
├── proto/                         # Shared protobuf contracts
│   └── stream.proto               # DSP stream service definition
│
├── .github/workflows/             # CI/CD pipelines
│   └── ci.yml                     # Build and test workflow
│
├── go.mod
└── README.md
```

## Quick Start

### Requirements
- Go 1.22+
- C++ Compiler (GCC/Clang)
- CMake 3.10+
- Protobuf Compiler (`protoc`)

### Build & Run

**1. Build the C++ DSP Library**
```bash
cd cpp
mkdir build && cd build
cmake ..
make
cd ../..
```

**2. Start the gRPC Server**
In a terminal, run the server:
```bash
go run ./cmd/server
```
*The server listens on `:50051`.*

**3. Generate Test Audio & Stream**
In a separate terminal, generate a test `test.wav` file (440Hz sine wave) and stream it through the gRPC server:
```bash
# Generate test.wav
go run ./cmd/genwav

# Stream through DSP engine -> output.wav
go run ./cmd/client test.wav output.wav

# Analyze the math (Zero-copy in-place mutation)
go run ./cmd/analyze
```

## API Contract (gRPC)

The communication relies on a bidirectional gRPC stream. The client sends a continuous flow of raw byte chunks, and the server returns the mutated byte chunks instantly.

```protobuf
service DSPStream {
  // Bidirectional streaming of audio chunks
  rpc ProcessAudio(stream AudioChunk) returns (stream AudioChunk);
}

message AudioChunk {
  // Raw bytes (reinterpreted as float32 in memory)
  bytes data = 1;
}
```

## Architecture

The project is designed to avoid memory leaks, allocations during streaming, and Go GC pauses.

```
                  ┌───────────────────────────────┐
                  │         Go gRPC Server        │
                  │                               │
[Client] ──────>  │  1. Receive gRPC []byte       │
 (Stream)         │  2. Cast to []float32         │
                  │     (unsafe.Pointer)          │
                  │  3. Pass pointer via CGO ──── │ ───┐
                  │                               │    │
[Client] <──────  │  5. Send mutated []byte       │    │
 (Stream)         │     back to client            │    │
                  └───────────────────────────────┘    │
                                                       │
                                                       ↓
                  ┌───────────────────────────────┐
                  │          C++ DSP Core         │
                  │                               │
                  │  4. Apply Biquad Filter       │
                  │     IN-PLACE mutation         │
                  │     (No allocations)          │
                  └───────────────────────────────┘
```

## Testing & CI/CD
The repository is equipped with a GitHub Actions workflow (`.github/workflows/ci.yml`) that validates code formatting, compiles the C++ static library, regenerates Protobuf contracts, and runs the Go test suite utilizing `bufconn` for in-memory gRPC testing.

```bash
go test -v ./...
```
