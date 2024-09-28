// в файле содержатся реализации методов gRPC-сервера для работы с получением списка директорий пользователя
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

// Folders реализует метод интерфейса UnimplementedGophKeeperServer. Получает список директорий
func (s *Server) Folders(ctx context.Context, r *pb.FoldersRequest) (*pb.FoldersResponse, error) {
	username, err := usernameFromToken(ctx, s.config.Secret())
	if err != nil {
		s.log.Error("запрос на сохранение данных", slog.String("ошибка", err.Error()))
		return nil, status.Error(codes.Unauthenticated, "токен не был передан или не прошел проверку")
	}
	resp, err := s.db.Folders(ctx, model.StorageFoldersRequest{User: username})
	if err != nil {
		s.log.Error("запрос папок у пользователя", slog.String("пользователь", username), slog.String("ошибка", err.Error()))
	}
	if errors.Is(err, storage.ErrKeepFolderNotExist) {
		return nil, status.Error(codes.NotFound, "ничего не найдено")
	}
	if err != nil {
		return nil, err
	}
	return &pb.FoldersResponse{
		Folders: resp.Folders,
	}, nil
}
