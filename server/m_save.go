package server

import (
	"context"
	"errors"
	"log/slog"

	pb "github.com/kTowkA/GophKeeper/grpc"
	"github.com/kTowkA/GophKeeper/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Save(ctx context.Context, r *pb.SaveRequest) (*pb.SaveResponse, error) {
	username, err := usernameFromToken(ctx, s.config.Secret())
	if err != nil {
		s.log.Error("запрос на сохранение данных", slog.String("ошибка", err.Error()))
		return nil, status.Error(codes.Unauthenticated, "токен не был передан или не прошел проверку")
	}
	_, err = s.db.Save(
		ctx,
		model.StorageSaveRequest{
			User:  username,
			Value: convertSaveRequestToModelKeeperElement(r),
		},
	)
	if err != nil {
		s.log.Error("сохранение данных", slog.String("пользователь", username), slog.String("ошибка", err.Error()))
		return nil, errors.New("произошла ошибка при сохранении данных")
	}
	return &pb.SaveResponse{
		SaveStatus:  true,
		SaveMessage: "ok",
	}, nil

}

func convertSaveRequestToModelKeeperElement(r *pb.SaveRequest) model.KeeperElement {
	mke := model.KeeperElement{
		Title:       r.Value.Title,
		Description: r.Value.Description,
		Type:        r.Value.Type,
		Values:      make([]model.KeeperValue, len(r.Value.Values)),
	}
	for i := range r.Value.Values {
		mke.Values[i] = model.KeeperValue{
			Title:       r.Value.Values[i].Title,
			Description: r.Value.Values[i].Description,
			Value:       r.Value.Values[i].Value,
		}
	}
	return mke
}
