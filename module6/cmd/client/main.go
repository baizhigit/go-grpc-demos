package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/baizhigit/go-grpc-demos/module6/internal/config"
	"github.com/baizhigit/go-grpc-demos/module6/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	cfg := config.Config{
		MethodConfig: []config.MethodConfig{{
			Name: []config.NameConfig{{
				Service: "config.ConfigService",
				Method:  "LongRunning",
			}},
			Timeout: "10s",
		}},
	}

	serviceConfig, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalf("failed to marshal config: %v", err)
	}

	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(string(serviceConfig)),
	)

	client := proto.NewConfigServiceClient(conn)

	_, err = client.LongRunning(ctx, &proto.LongRunningRequest{})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("successful response received")
}
