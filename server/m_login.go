package server

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v4"
	pb "github.com/kTowkA/GophKeeper/grpc"
	"github.com/kTowkA/GophKeeper/internal/model"
	"github.com/kTowkA/GophKeeper/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const TokenExp12Hours = 12 * time.Hour

func (s *Server) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp, err := s.db.PasswordHash(
		ctx,
		model.StoragePasswordHashRequest{
			Login: r.Login,
		},
	)
	switch {
	case errors.Is(err, storage.ErrUserIsNotExist):
		s.log.Debug("попытка входа несуществующего пользователя", slog.String("логин", r.Login))
		return nil, status.Error(codes.NotFound, "пользователь не найден")
	case err != nil:
		s.log.Error("запрос хеша пароля у пользователя", slog.String("логин", r.Login), slog.String("ошибка", err.Error()))
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(resp.PasswordHash), []byte(r.Password))
	if err != nil {
		s.log.Debug("попытка входа с неверным паролем", slog.String("логин", r.Login))
		return nil, status.Error(codes.PermissionDenied, "пароль неверен")
	}
	token, err := generateToken(r.Login, s.config.Secret())
	if err != nil {
		s.log.Error("генерация токена", slog.String("логин", r.Login), slog.String("ошибка", err.Error()))
		return nil, errors.New("ошибка генерации токена")
	}
	return &pb.LoginResponse{
		LoginStatus:  true,
		LoginMessage: "ok",
		Token:        token,
	}, nil
}

type Claims struct {
	jwt.RegisteredClaims
	User string
}

func generateToken(user string, secret string) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp12Hours)),
		},
		User: user,
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
