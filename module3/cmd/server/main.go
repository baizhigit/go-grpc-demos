package main

import (
	"log"
	"net"

	"github.com/baizhigit/go-grpc-demos/module3/internal/stream"
	"github.com/baizhigit/go-grpc-demos/module3/proto"
	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()

	streamingService := &stream.Service{}

	proto.RegisterFileUploadServiceServer(grpcServer, streamingService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("starting grpc server on address: %s", ":50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
