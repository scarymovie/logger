package main

import (
	"context"

	pb "github.com/scarymovie/logger/cmd/grpc-server/proto"
	"github.com/scarymovie/logger/slogx"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (s *greeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	slogx.Info(ctx, "handling SayHello", slogx.String("name", req.Name))
	return &pb.HelloReply{Message: "Hello, " + req.Name + "!"}, nil
}
