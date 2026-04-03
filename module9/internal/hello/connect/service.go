package connect

import (
	"context"
	"errors"
	"fmt"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"github.com/baizhigit/go-grpc-demos/module9/proto"
	"github.com/baizhigit/go-grpc-demos/module9/proto/protoconnect"
	googleproto "google.golang.org/protobuf/proto"
)

type (
	RequestValidator interface {
		Validate(msg googleproto.Message, opts ...protovalidate.ValidationOption) error
	}

	service struct {
		protoconnect.UnimplementedHelloServiceHandler
		validator RequestValidator
	}
)

func NewService(validator RequestValidator) (*service, error) {
	if validator == nil {
		return nil, errors.New("validator cannot be nil")
	}

	return &service{
		validator: validator,
	}, nil
}

func (s *service) SayHello(ctx context.Context, c *connect.Request[proto.SayHelloRequest]) (*connect.Response[proto.SayHelloResponse], error) {
	if err := s.validator.Validate(c.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	return connect.NewResponse(&proto.SayHelloResponse{
		Message: fmt.Sprintf("Hello %s", c.Msg.GetName()),
	}), nil
}
