package app

import (
	"context"
	"fmt"
	"runtime/debug"

	grpchandler "github.com/nghiatrann0502/trinity/internal/video/adapters/grpcHandler"
	"github.com/nghiatrann0502/trinity/internal/video/core/ports"
	"github.com/nghiatrann0502/trinity/pkg/logger"
	"github.com/nghiatrann0502/trinity/proto/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func recoveryInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := debug.Stack()
				log.Error("GRPC internal error", fmt.Errorf("%v", r), map[string]interface{}{
					"stack_trace": stackTrace,
				})

				err = status.Errorf(codes.Internal, "Internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

func NewGRPCServer(log logger.Logger, svc ports.Service) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(recoveryInterceptor(log)),
	)
	userGRPCServer := grpchandler.NewGRPCHandler(log, svc)
	proto.RegisterVideoServiceServer(grpcServer, userGRPCServer)

	return grpcServer
}
