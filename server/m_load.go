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

func (s *Server) Load(ctx context.Context, r *pb.LoadRequest) (*pb.LoadResponse, error) {
	username, err := usernameFromToken(ctx, s.config.Secret())
	if err != nil {
		s.log.Error("запрос на сохранение данных", slog.String("ошибка", err.Error()))
		return nil, status.Error(codes.Unauthenticated, "токен не был передан или не прошел проверку")
	}
	resp, err := s.db.Load(ctx, model.StorageLoadRequest{
		User:               username,
		TitleKeeperElement: r.Title,
	})
	if err != nil {
		s.log.Error("загрузка данных", slog.String("пользователь", username), slog.String("ошибка", err.Error()))
		return nil, errors.New("произошла ошибка при запросе данных")
	}
	return convertModelStorageLoadResponseToLoadResponse(resp), nil
}

func convertModelStorageLoadResponseToLoadResponse(r model.StorageLoadResponse) *pb.LoadResponse {
	lr := &pb.LoadResponse{
		Value: &pb.KeeperElement{
			Title:       r.TitleKeeperElement.Title,
			Description: r.TitleKeeperElement.Description,
			Type:        r.TitleKeeperElement.Type,
			Values:      make([]*pb.KeeperElement_KeeperValue, len(r.TitleKeeperElement.Values)),
		},
	}
	for i := range r.TitleKeeperElement.Values {
		lr.Value.Values[i] = &pb.KeeperElement_KeeperValue{
			Title:       r.TitleKeeperElement.Values[i].Title,
			Description: r.TitleKeeperElement.Values[i].Description,
			Value:       r.TitleKeeperElement.Values[i].Value,
		}
	}
	return lr
}
