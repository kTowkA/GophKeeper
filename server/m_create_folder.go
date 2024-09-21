package server

import (
	"context"
	"errors"
	"log/slog"

	pb "github.com/kTowkA/GophKeeper/grpc"
	"github.com/kTowkA/GophKeeper/internal/model"
	"github.com/kTowkA/GophKeeper/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateFolder(ctx context.Context, r *pb.CreateFolderRequest) (*pb.CreateFolderResponse, error) {
	username, err := usernameFromToken(ctx, s.config.Secret())
	if err != nil {
		s.log.Error("запрос на создание папки", slog.String("ошибка", err.Error()))
		return nil, status.Error(codes.Unauthenticated, "токен не был передан или не прошел проверку")
	}
	resp, err := s.db.CreateFolder(
		ctx,
		model.StorageCreateFolderRequest{
			User:        username,
			Folder:      r.Title,
			Description: r.Description,
		},
	)
	if err != nil {
		s.log.Error("запрос на создание папки", slog.String("пользователь", username), slog.String("папка", r.Title), slog.String("ошибка", err.Error()))
	}
	if errors.Is(err, storage.ErrKeepFolderIsExist) {
		return nil, status.Error(codes.AlreadyExists, "папка уже существует")
	}
	if err != nil {
		return &pb.CreateFolderResponse{}, err
	}
	return &pb.CreateFolderResponse{
		CreateFolderStatus:  true,
		CreateFolderMessage: "ok",
		CreateFolderUuid:    resp.FolderID.String(),
	}, nil
}
