package server

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	pb "github.com/kTowkA/GophKeeper/grpc"
	"github.com/kTowkA/GophKeeper/internal/model"
	"github.com/kTowkA/GophKeeper/internal/storage"
	mocks "github.com/kTowkA/GophKeeper/internal/storage/mocs"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerTest struct {
	suite.Suite
	ms *mocks.Storager
	gs *Server
}

type Test struct {
	name          string
	request       any
	ctx           context.Context
	wantError     bool
	wantgRPCError codes.Code
	wantResponse  any
	mock          func()
}

func (suite *ServerTest) SetupSuite() {
	suite.Suite.T().Log("Suite setup")

	suite.ms = new(mocks.Storager)

	suite.gs = new(Server)
	suite.gs.db = suite.ms
	suite.gs.log = slog.Default()
}

func (suite *ServerTest) TearDownSuite() {
	//заканчиваем работу с тестовым сценарием
	defer suite.Suite.T().Log("Suite test is ended")

	// смотрим чтобы все описанные вызовы были использованы
	suite.ms.AssertExpectations(suite.T())
}

func (suite *ServerTest) TestRegister() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tests := []Test{
		{
			name: "пользователь существует",
			request: &pb.RegisterRequest{
				Login:    "1",
				Password: "2111qweQASXC!",
			},
			ctx:           ctx,
			wantgRPCError: codes.AlreadyExists,
			wantError:     true,
			mock: func() {
				suite.ms.On("Register", mock.Anything, mock.AnythingOfType("model.StorageRegisterRequest")).Return(model.StorageRegisterResponse{}, storage.ErrLoginIsAlreadyOccupied).Once()
			},
		},
		{
			name: "ошибка при сохранении в базу данных",
			request: &pb.RegisterRequest{
				Login:    "1",
				Password: "2111qweQASXC!",
			},
			ctx:       ctx,
			wantError: true,
			mock: func() {
				suite.ms.On("Register", mock.Anything, mock.AnythingOfType("model.StorageRegisterRequest")).Return(model.StorageRegisterResponse{}, errors.New("ошибка в базе данных")).Once()
			},
		},
		{
			name: "успешная регистрация",
			request: &pb.RegisterRequest{
				Login:    "1",
				Password: "2111qweQASXC!",
			},
			ctx:       ctx,
			wantError: false,
			wantResponse: &pb.RegisterResponse{
				RegisterStatus:  true,
				RegisterMessage: "ok",
			},
			mock: func() {
				suite.ms.On("Register", mock.Anything, mock.AnythingOfType("model.StorageRegisterRequest")).Return(model.StorageRegisterResponse{}, nil).Once()
			},
		},
	}

	for _, t := range tests {
		if t.mock != nil {
			t.mock()
		}
		resp, err := suite.gs.Register(t.ctx, (t.request).(*pb.RegisterRequest))
		if !t.wantError {
			suite.NoError(err, t.name)
			suite.EqualValues((t.wantResponse).(*pb.RegisterResponse), resp)
			continue
		}
		if t.wantgRPCError > 0 {
			if e, ok := status.FromError(err); ok {
				suite.EqualValues(t.wantgRPCError, e.Code(), t.name)
			} else {
				suite.Fail("должна содержаться ошибка", t.name)
			}
			continue
		}
		suite.Error(err, t.name)
	}
}

func (suite *ServerTest) TestLogin() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tests := []Test{
		{
			name: "пользователь не существует",
			request: &pb.LoginRequest{
				Login:    "1",
				Password: "2",
			},
			ctx:           ctx,
			wantError:     true,
			wantgRPCError: codes.NotFound,
			mock: func() {
				suite.ms.On("PasswordHash", mock.Anything, model.StoragePasswordHashRequest{Login: "1"}).Return(model.StoragePasswordHashResponse{}, storage.ErrLoginIsNotExist).Once()
			},
		},
		{
			name: "ошибка в БД",
			request: &pb.LoginRequest{
				Login:    "1",
				Password: "2",
			},
			ctx:       ctx,
			wantError: true,
			mock: func() {
				suite.ms.On("PasswordHash", mock.Anything, model.StoragePasswordHashRequest{Login: "1"}).Return(model.StoragePasswordHashResponse{}, errors.New("ошибка в базе данных")).Once()
			},
		},
		{
			name: "неверный пароль",
			request: &pb.LoginRequest{
				Login:    "1",
				Password: "2",
			},
			ctx:           ctx,
			wantError:     true,
			wantgRPCError: codes.PermissionDenied,
			mock: func() {
				suite.ms.On("PasswordHash", mock.Anything, model.StoragePasswordHashRequest{Login: "1"}).Return(model.StoragePasswordHashResponse{PasswordHash: "123"}, nil).Once()
			},
		},
		{
			name: "все хорошо",
			request: &pb.LoginRequest{
				Login:    "1",
				Password: "2",
			},
			ctx:       ctx,
			wantError: false,
			mock: func() {
				pwd, _ := bcrypt.GenerateFromPassword([]byte("2"), bcrypt.DefaultCost)
				suite.ms.On("PasswordHash", mock.Anything, model.StoragePasswordHashRequest{Login: "1"}).Return(model.StoragePasswordHashResponse{PasswordHash: string(pwd)}, nil).Once()
			},
		},
	}

	for _, t := range tests {
		if t.mock != nil {
			t.mock()
		}
		_, err := suite.gs.Login(t.ctx, (t.request).(*pb.LoginRequest))
		if !t.wantError {
			suite.NoError(err, t.name)
			continue
		}
		if t.wantgRPCError > 0 {
			if e, ok := status.FromError(err); ok {
				suite.EqualValues(t.wantgRPCError, e.Code(), t.name)
			} else {
				suite.Fail("должна содержаться ошибка", t.name)
			}
			continue
		}
		suite.Error(err, t.name)
	}
}
func TestAppSuite(t *testing.T) {
	suite.Run(t, new(ServerTest))
}
