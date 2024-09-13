package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/golang-jwt/jwt/v4"
	pb "github.com/kTowkA/GophKeeper/grpc"
	"github.com/kTowkA/GophKeeper/internal/config"
	"github.com/kTowkA/GophKeeper/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	TokenTitle = "token"
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

func usernameFromToken(ctx context.Context, secret string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "не найден контекст для проверки доступа")
	}
	tokenCtx := md.Get(TokenTitle)
	if len(tokenCtx) == 0 {
		return "", status.Error(codes.Unauthenticated, "не найден токен в переданном контексте")
	}
	tokenString := tokenCtx[0]

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неожиданный метод подписи: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("токен не прошел проверку")
	}
	return claims.User, nil
}
