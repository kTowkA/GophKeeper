package server

import (
	"context"
	"errors"
	"log/slog"

	pb "github.com/kTowkA/GophKeeper/grpc"
	"github.com/kTowkA/GophKeeper/internal/model"
	"github.com/kTowkA/GophKeeper/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Register(ctx context.Context, r *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("получение хеша пароля", slog.String("логин", r.Login), slog.String("ошибка", err.Error()))
		return nil, err
	}
	_, err = s.db.Register(
		ctx,
		model.StorageRegisterRequest{
			Login:    r.Login,
			Password: string(hash),
		},
	)
	switch {
	case errors.Is(err, storage.ErrLoginIsAlreadyOccupied):
		s.log.Debug("попытка регистрации с существующим логином", slog.String("логин", r.Login))
		return nil, status.Error(codes.AlreadyExists, err.Error())
	case err != nil:
		s.log.Error("регистрация пользователя", slog.String("логин", r.Login), slog.String("ошибка", err.Error()))
		return nil, err
	}
	return &pb.RegisterResponse{
		RegisterStatus:  true,
		RegisterMessage: "ok",
	}, nil
}
