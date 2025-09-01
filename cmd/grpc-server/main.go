package main

import (
	"context"
	"log/slog"
	"net"

	"github.com/google/uuid"
	pb "github.com/scarymovie/logger/cmd/grpc-server/proto"
	"github.com/scarymovie/logger/slogx"
	"google.golang.org/grpc"
)

func loggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	traceID := uuid.New().String()
	ctx = slogx.WithContext(ctx, slogx.String("trace_id", traceID))

	slogx.Info(ctx, "incoming gRPC request", slogx.String("method", info.FullMethod))

	resp, err := handler(ctx, req)
	if err != nil {
		slogx.Error(ctx, "gRPC request failed", slogx.String("error", err.Error()))
		return nil, err
	}

	slogx.Info(ctx, "gRPC request completed")
	return resp, nil
}

func main() {
	slogx.MustConfigure(slogx.Config{
		Format:       "json",
		Level:        slog.LevelDebug,
		AddSource:    true,
		DefaultAttrs: []slog.Attr{slogx.String("service", "grpc-greeter")},
	})

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	pb.RegisterGreeterServer(grpcServer, &greeterServer{})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		slogx.Error(context.Background(), "failed to listen", slogx.String("error", err.Error()))
		return
	}

	slogx.Info(context.Background(), "gRPC server starting", slogx.String("addr", ":50051"))
	if err := grpcServer.Serve(lis); err != nil {
		slogx.Error(context.Background(), "server failed", slogx.String("error", err.Error()))
	}
}
