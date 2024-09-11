package server

import (
	"context"
	"log/slog"

	pb "github.com/kTowkA/GophKeeper/grpc"
	"github.com/kTowkA/GophKeeper/internal/config"
	"github.com/kTowkA/GophKeeper/internal/storage"
)

type Server struct {
	log    *slog.Logger
	config config.ConfigServer
	db     storage.Storager
	pb.UnimplementedGophKeeperServer
}

func NewServer() (*Server, error) {
	return nil, nil
}

func (s *Server) Run(ctx context.Context) error {
	return nil
}
