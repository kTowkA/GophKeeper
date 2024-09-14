package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	pb "github.com/kTowkA/GophKeeper/grpc"
	"github.com/kTowkA/GophKeeper/internal/storage"
	"github.com/kTowkA/GophKeeper/server/config"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	TokenTitle = "token"
)

type Server struct {
	log    *slog.Logger
	config config.Config
	db     storage.Storager
	pb.UnimplementedGophKeeperServer
}

func newServer(config config.Config, db storage.Storager) (*Server, error) {
	if db == nil {
		return nil, errors.New("для корректной работы сервера необходимо установить хранилище")
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: config.LogLevel()}))
	server := Server{
		log:    log,
		config: config,
		db:     db,
	}
	return &server, nil
}

func Run(ctx context.Context, db storage.Storager, config config.Config) error {
	s, err := newServer(config, db)
	if err != nil {
		return err
	}
	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(),
	))

	gr, grCtx := errgroup.WithContext(ctx)

	gr.Go(func() error {
		defer s.log.Info("сервер был остановлен")

		<-grCtx.Done()

		gRPCServer.GracefulStop()

		return nil
	})

	gr.Go(func() error {
		pb.RegisterGophKeeperServer(gRPCServer, s)

		l, err := net.Listen("tcp", s.config.Address())
		if err != nil {
			return err
		}

		s.log.Info("запуск сервера", slog.String("адрес", s.config.Address()))

		if err := gRPCServer.Serve(l); err != nil {
			return err
		}
		return nil
	})
	return gr.Wait()
}

func usernameFromToken(ctx context.Context, secret string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("не найден контекст для проверки доступа")
	}
	tokenCtx := md.Get(TokenTitle)
	if len(tokenCtx) == 0 {
		return "", errors.New("не найден токен в переданном контексте")
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
		return "", errors.New("токен не прошел проверку")
	}
	return claims.User, nil
}
