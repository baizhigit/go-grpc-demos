package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/baizhigit/go-grpc-demos/module5/internal/token"
	"github.com/baizhigit/go-grpc-demos/module5/proto"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	if err := run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	// get jwt secret from env var
	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return errors.New("JWT_SECRET must be provided")
	}

	// initialise our middleware
	middleware, err := token.NewMiddleware([]byte(jwtSecret))
	if err != nil {
		return fmt.Errorf("failed to initialise middleware: %w", err)
	}

	// create our server
	server := grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryAuthMiddleware))

	// initialize our handler
	tokenService := token.Service{}

	// register our service
	proto.RegisterTokenServiceServer(server, tokenService)

	// create an errgroup & start our server
	g, ctx := errgroup.WithContext(ctx)

	const addr = ":50051"

	g.Go(func() error {
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to listen on address %q: %w", addr, err)
		}

		log.Printf("server started on address %q", addr)

		if err := server.Serve(lis); err != nil {
			return fmt.Errorf("failed to serve grpc service: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		<-ctx.Done()

		server.GracefulStop()

		return nil
	})

	return g.Wait()
}
