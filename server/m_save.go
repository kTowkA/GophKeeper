// в файле содержатся реализации методов gRPC-сервера для работы с сохранением конкретных данных
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

// Save реализует метод интерфейса UnimplementedGophKeeperServer. Сохранение данных
func (s *Server) Save(ctx context.Context, r *pb.SaveRequest) (*pb.SaveResponse, error) {
	username, err := usernameFromToken(ctx, s.config.Secret())
	if err != nil {
		s.log.Error("запрос на сохранение данных", slog.String("ошибка", err.Error()))
		return nil, status.Error(codes.Unauthenticated, "токен не был передан или не прошел проверку")
	}
	resp, err := s.db.Save(
		ctx,
		model.StorageSaveRequest{
			User:   username,
			Folder: r.Folder,
			Value:  convertSaveRequestToModelKeeperElement(r),
		},
	)
	if err != nil {
		s.log.Error("сохранение данных", slog.String("пользователь", username), slog.String("название", r.Value.Title), slog.String("ошибка", err.Error()))
	}
	if errors.Is(err, storage.ErrKeepValueIsExist) {
		return nil, status.Error(codes.AlreadyExists, "данные под таким названием уже существуют")
	}
	if err != nil {
		return nil, errors.New("произошла ошибка при сохранении данных")
	}
	return &pb.SaveResponse{
		SaveStatus:    true,
		SaveMessage:   "ok",
		SaveValueUuid: resp.ValueID.String(),
	}, nil

}

func convertSaveRequestToModelKeeperElement(r *pb.SaveRequest) model.KeeperValue {
	mke := model.KeeperValue{
		Title:       r.Value.Title,
		Description: r.Value.Description,
		Value:       r.Value.Value,
	}
	return mke
}
