# gRPC with Go

This project demonstrates how to build and use gRPC services in Go.

## Features

- Define gRPC services using Protocol Buffers
- Implement server and client in Go
- Example request/response handling

## Prerequisites

- Go 1.18+
- `protoc` Protocol Buffers compiler
- `protoc-gen-go` and `protoc-gen-go-grpc` plugins

## Setup

1. Clone the repository:

   ```sh
   git clone https://github.com/MagicRodri/grpc_with_go.git
   cd grpc_with_go
   ```

2. Install dependencies:

   ```sh
   go mod tidy
   ```

3. Generate gRPC code:
   ```sh
   make proto
   ```

## Running

- Start the server:

  ```sh
  go run server/main.go
  ```
