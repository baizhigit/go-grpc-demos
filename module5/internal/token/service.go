package token

import (
	"context"

	"github.com/baizhigit/go-grpc-demos/module5/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	proto.UnimplementedTokenServiceServer
}

func (s Service) Validate(ctx context.Context, _ *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	claims, ok := ctx.Value(claimsKey).(map[string]string)
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "claims missing from context")
	}

	return &proto.ValidateResponse{Claims: claims}, nil
}
