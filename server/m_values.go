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

func (s *Server) Values(ctx context.Context, r *pb.ValuesInFolderRequest) (*pb.ValuesInFolderResponse, error) {
	username, err := usernameFromToken(ctx, s.config.Secret())
	if err != nil {
		s.log.Error("запрос на сохранение данных", slog.String("ошибка", err.Error()))
		return nil, status.Error(codes.Unauthenticated, "токен не был передан или не прошел проверку")
	}
	resp, err := s.db.Values(ctx, model.StorageValuesRequest{
		User:   username,
		Folder: r.Folder,
	})
	if err != nil {
		s.log.Error("запрос заголовков данных у пользователя", slog.String("пользователь", username), slog.String("папка", r.Folder), slog.String("ошибка", err.Error()))
	}
	if errors.Is(err, storage.ErrKeepFolderNotExist) {
		return nil, status.Error(codes.NotFound, "ничего не найдено")
	}
	if errors.Is(err, storage.ErrKeepValueNotExist) {
		return nil, status.Error(codes.NotFound, "ничего не найдено")
	}
	if err != nil {
		return nil, err
	}
	return &pb.ValuesInFolderResponse{
		Values: resp.Values,
	}, nil
}
