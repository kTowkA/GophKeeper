package server

import (
	"context"

	pb "github.com/kTowkA/GophKeeper/grpc"
)

func (s *Server) Save(ctx context.Context, r *pb.SaveRequest) (*pb.SaveResponse, error) {
	return nil, nil

}
