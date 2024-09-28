// в файле содержатся реализации методов gRPC-сервера для работы с получением конкретных данных
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

// Load реализует метод интерфейса UnimplementedGophKeeperServer. Получает конкретные данные пользователя в конкретной директории
func (s *Server) Load(ctx context.Context, r *pb.LoadRequest) (*pb.LoadResponse, error) {
	username, err := usernameFromToken(ctx, s.config.Secret())
	if err != nil {
		s.log.Error("запрос на сохранение данных", slog.String("ошибка", err.Error()))
		return nil, status.Error(codes.Unauthenticated, "токен не был передан или не прошел проверку")
	}
	resp, err := s.db.Load(ctx, model.StorageLoadRequest{
		User:   username,
		Folder: r.Folder,
		Title:  r.Title,
	})
	if err != nil {
		s.log.Error("загрузка данных", slog.String("пользователь", username), slog.String("ошибка", err.Error()))
	}
	if errors.Is(err, storage.ErrKeepValueNotExist) {
		return nil, status.Error(codes.NotFound, "данные не найдены")
	}
	if err != nil {
		return nil, errors.New("произошла ошибка при запросе данных")
	}
	return convertModelStorageLoadResponseToLoadResponse(resp), nil
}

func convertModelStorageLoadResponseToLoadResponse(r model.StorageLoadResponse) *pb.LoadResponse {
	lr := &pb.LoadResponse{
		Value: &pb.KeeperValue{
			Title:       r.Value.Title,
			Description: r.Value.Description,
			Value:       r.Value.Value,
		},
	}
	return lr
}
