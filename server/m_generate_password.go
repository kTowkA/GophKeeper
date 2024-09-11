package server

import (
	"context"

	pb "github.com/kTowkA/GophKeeper/grpc"
)

func (s *Server) GeneratePassword(context.Context, *pb.GeneratePasswordRequest) (*pb.GeneratePasswordResponse, error) {
	return nil, nil
}
