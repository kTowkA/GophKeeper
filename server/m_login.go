package server

import (
	"context"

	pb "github.com/kTowkA/GophKeeper/grpc"
)

func (s *Server) Login(context.Context, *pb.LoginRequest) (*pb.LoginResponse, error) {
	return nil, nil
}
