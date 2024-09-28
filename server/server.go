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
	// tokenTitle ключ токена в контексте
	tokenTitle = "token"
)

// Server стуктура сервера. При вызове метода Run создается автоматически
// реализует интерфейс UnimplementedGophKeeperServer
type Server struct {
	log    *slog.Logger
	config config.Config
	db     storage.Storager
	pb.UnimplementedGophKeeperServer
}

// newServer создание экземпляра сервера
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

// Run запуск сервера по адресу указанному в config. При завершении работы по сигналу контекста ошибку не возвращает
func Run(ctx context.Context, db storage.Storager, config config.Config) error {
	s, err := newServer(config, db)
	if err != nil {
		return err
	}
	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(),
	))

	gr, grCtx := errgroup.WithContext(ctx)

	// горутина завершения работы
	gr.Go(func() error {
		defer s.log.Info("сервер был остановлен")

		<-grCtx.Done()

		gRPCServer.GracefulStop()

		return nil
	})

	// горутина запуска gRPC сервера
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

// usernameFromToken получение имени пользователя из переданного токена (+ проверка валидации токена)
// нужна так как в запросах к хранилищу указывается пользователь
func usernameFromToken(ctx context.Context, secret string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("не найден контекст для проверки доступа")
	}
	tokenCtx := md.Get(tokenTitle)
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
