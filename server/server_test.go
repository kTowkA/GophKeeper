package server

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	pb "github.com/kTowkA/GophKeeper/grpc"
	"github.com/kTowkA/GophKeeper/internal/model"
	"github.com/kTowkA/GophKeeper/internal/storage"
	mocks "github.com/kTowkA/GophKeeper/internal/storage/mocs"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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
				suite.ms.On("PasswordHash", mock.Anything, model.StoragePasswordHashRequest{Login: "1"}).Return(model.StoragePasswordHashResponse{}, storage.ErrUserIsNotExist).Once()
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

func (suite *ServerTest) TestGeneratePassword() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := suite.gs.GeneratePassword(ctx, &pb.GeneratePasswordRequest{Length: 15})
	suite.NoError(err)
	suite.EqualValues(15, utf8.RuneCountInString(resp.Password))
	suite.NoError(validatePassword(resp.Password))
}

func (suite *ServerTest) TestCreateFolder() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := "user"
	secret := ""
	token, err := generateToken(username, secret)
	suite.NoError(err)
	validCtx := metadata.NewIncomingContext(ctx, metadata.New(map[string]string{TokenTitle: token}))

	tests := []Test{
		{
			name: "пользователь не существует",
			request: &pb.CreateFolderRequest{
				Title: "1",
			},
			ctx:           ctx,
			wantError:     true,
			wantgRPCError: codes.Unauthenticated,
		},
		{
			name: "папка уже существует",
			request: &pb.CreateFolderRequest{
				Title: "1",
			},
			ctx:           validCtx,
			wantError:     true,
			wantgRPCError: codes.AlreadyExists,
			mock: func() {
				suite.ms.On("CreateFolder", mock.Anything, model.StorageCreateFolderRequest{User: username, Folder: "1"}).Return(model.StorageCreateFolderResponse{}, storage.ErrKeepFolderIsExist).Once()
			},
		},
		{
			name: "ошибка в БД",
			request: &pb.CreateFolderRequest{
				Title: "1",
			},
			ctx:       validCtx,
			wantError: true,
			mock: func() {
				suite.ms.On("CreateFolder", mock.Anything, model.StorageCreateFolderRequest{User: username, Folder: "1"}).Return(model.StorageCreateFolderResponse{}, errors.New("ошибка при создании папки")).Once()
			},
		},
		{
			name: "все ок",
			request: &pb.CreateFolderRequest{
				Title: "1",
			},
			ctx:       validCtx,
			wantError: false,
			mock: func() {
				suite.ms.On("CreateFolder", mock.Anything, model.StorageCreateFolderRequest{User: username, Folder: "1"}).Return(model.StorageCreateFolderResponse{FolderID: uuid.New()}, nil).Once()
			},
		},
	}
	for _, t := range tests {
		if t.mock != nil {
			t.mock()
		}
		_, err := suite.gs.CreateFolder(t.ctx, (t.request).(*pb.CreateFolderRequest))
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
func (suite *ServerTest) TestFolders() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := "user"
	secret := ""
	token, err := generateToken(username, secret)
	suite.NoError(err)
	validCtx := metadata.NewIncomingContext(ctx, metadata.New(map[string]string{TokenTitle: token}))

	tests := []Test{
		{
			name:          "пользователь не существует",
			request:       &pb.FoldersRequest{},
			ctx:           ctx,
			wantError:     true,
			wantgRPCError: codes.Unauthenticated,
		},
		{
			name:          "папок нет",
			request:       &pb.FoldersRequest{},
			ctx:           validCtx,
			wantError:     true,
			wantgRPCError: codes.NotFound,
			mock: func() {
				suite.ms.On("Folders", mock.Anything, model.StorageFoldersRequest{User: username}).Return(model.StorageFoldersResponse{}, storage.ErrKeepFolderNotExist).Once()
			},
		},
		{
			name:      "ошибка в БД",
			request:   &pb.FoldersRequest{},
			ctx:       validCtx,
			wantError: true,
			mock: func() {
				suite.ms.On("Folders", mock.Anything, model.StorageFoldersRequest{User: username}).Return(model.StorageFoldersResponse{}, errors.New("ошибка при запросе папок")).Once()
			},
		},
		{
			name:      "все ок",
			request:   &pb.FoldersRequest{},
			ctx:       validCtx,
			wantError: false,
			mock: func() {
				suite.ms.On("Folders", mock.Anything, model.StorageFoldersRequest{User: username}).Return(model.StorageFoldersResponse{}, nil).Once()
			},
		},
	}
	for _, t := range tests {
		if t.mock != nil {
			t.mock()
		}
		_, err := suite.gs.Folders(t.ctx, (t.request).(*pb.FoldersRequest))
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
func (suite *ServerTest) TestValues() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := "user"
	secret := ""
	token, err := generateToken(username, secret)
	suite.NoError(err)
	validCtx := metadata.NewIncomingContext(ctx, metadata.New(map[string]string{TokenTitle: token}))

	tests := []Test{
		{
			name:          "пользователь не существует",
			request:       &pb.ValuesInFolderRequest{},
			ctx:           ctx,
			wantError:     true,
			wantgRPCError: codes.Unauthenticated,
		},
		{
			name: "папок нет",
			request: &pb.ValuesInFolderRequest{
				Folder: "folder",
			},
			ctx:           validCtx,
			wantError:     true,
			wantgRPCError: codes.NotFound,
			mock: func() {
				suite.ms.On("Values", mock.Anything, model.StorageValuesRequest{User: username, Folder: "folder"}).Return(model.StorageValuesResponse{}, storage.ErrKeepFolderNotExist).Once()
			},
		},
		{
			name: "значений нет",
			request: &pb.ValuesInFolderRequest{
				Folder: "folder",
			},
			ctx:           validCtx,
			wantError:     true,
			wantgRPCError: codes.NotFound,
			mock: func() {
				suite.ms.On("Values", mock.Anything, model.StorageValuesRequest{User: username, Folder: "folder"}).Return(model.StorageValuesResponse{}, storage.ErrKeepValueNotExist).Once()
			},
		},
		{
			name: "ошибка в БД",
			request: &pb.ValuesInFolderRequest{
				Folder: "folder",
			},
			ctx:       validCtx,
			wantError: true,
			mock: func() {
				suite.ms.On("Values", mock.Anything, model.StorageValuesRequest{User: username, Folder: "folder"}).Return(model.StorageValuesResponse{}, errors.New("ошибка при запросе содержимого папки")).Once()
			},
		},
		{
			name: "все хорошо",
			request: &pb.ValuesInFolderRequest{
				Folder: "folder",
			},
			ctx:       validCtx,
			wantError: false,
			mock: func() {
				suite.ms.On("Values", mock.Anything, model.StorageValuesRequest{User: username, Folder: "folder"}).Return(model.StorageValuesResponse{}, nil).Once()
			},
		},
	}
	for _, t := range tests {
		if t.mock != nil {
			t.mock()
		}
		_, err := suite.gs.Values(t.ctx, (t.request).(*pb.ValuesInFolderRequest))
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
func (suite *ServerTest) TestSave() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := "user"
	secret := ""
	token, err := generateToken(username, secret)
	suite.NoError(err)
	validCtx := metadata.NewIncomingContext(ctx, metadata.New(map[string]string{TokenTitle: token}))
	tests := []Test{
		{
			name:          "ошибка получения пользователя",
			ctx:           ctx,
			wantError:     true,
			wantgRPCError: codes.Unauthenticated,
			request: &pb.SaveRequest{
				Value: &pb.KeeperValue{},
			},
		},
		{
			name:      "ошибка в БД",
			ctx:       validCtx,
			wantError: true,
			request: &pb.SaveRequest{
				Value: &pb.KeeperValue{},
			},
			mock: func() {
				suite.ms.On("Save", mock.Anything, mock.AnythingOfType("model.StorageSaveRequest")).Return(model.StorageSaveResponse{}, errors.New("ошибка при сохранении данных")).Once()
			},
		},
		{
			name: "все хорошо",
			request: &pb.SaveRequest{
				Value: &pb.KeeperValue{},
			},
			ctx:       validCtx,
			wantError: false,
			mock: func() {
				suite.ms.On("Save", mock.Anything, mock.AnythingOfType("model.StorageSaveRequest")).Return(model.StorageSaveResponse{}, nil).Once()
			},
		},
		{
			name: "дубль",
			request: &pb.SaveRequest{
				Value: &pb.KeeperValue{},
			},
			ctx:           validCtx,
			wantError:     true,
			wantgRPCError: codes.AlreadyExists,
			mock: func() {
				suite.ms.On("Save", mock.Anything, mock.AnythingOfType("model.StorageSaveRequest")).Return(model.StorageSaveResponse{}, storage.ErrKeepValueIsExist).Once()
			},
		},
	}
	for _, t := range tests {
		if t.mock != nil {
			t.mock()
		}
		_, err := suite.gs.Save(t.ctx, (t.request).(*pb.SaveRequest))
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
func (suite *ServerTest) TestLoad() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := "user"
	secret := ""
	token, err := generateToken(username, secret)
	suite.NoError(err)
	validCtx := metadata.NewIncomingContext(ctx, metadata.New(map[string]string{TokenTitle: token}))
	tests := []Test{
		{
			name:          "ошибка получения пользователя",
			ctx:           ctx,
			wantError:     true,
			wantgRPCError: codes.Unauthenticated,
			request:       &pb.LoadRequest{},
		},
		{
			name:      "ошибка в БД",
			ctx:       validCtx,
			wantError: true,
			request: &pb.LoadRequest{
				Title:  "title",
				Folder: "folder",
			},
			mock: func() {
				suite.ms.On("Load", mock.Anything, model.StorageLoadRequest{
					User:   username,
					Folder: "folder",
					Title:  "title",
				}).Return(model.StorageLoadResponse{}, errors.New("ошибка при запросе данных")).Once()
			},
		},
		{
			name:      "все хорошо",
			ctx:       validCtx,
			wantError: false,
			request: &pb.LoadRequest{
				Title:  "title",
				Folder: "folder",
			},
			mock: func() {
				suite.ms.On("Load", mock.Anything, model.StorageLoadRequest{
					User:   username,
					Folder: "folder",
					Title:  "title",
				}).Return(model.StorageLoadResponse{}, nil).Once()
			},
		},
		{
			name:          "данных нет",
			ctx:           validCtx,
			wantError:     true,
			wantgRPCError: codes.NotFound,
			request: &pb.LoadRequest{
				Title:  "title",
				Folder: "folder",
			},
			mock: func() {
				suite.ms.On("Load", mock.Anything, model.StorageLoadRequest{
					User:   username,
					Folder: "folder",
					Title:  "title",
				}).Return(model.StorageLoadResponse{}, storage.ErrKeepValueNotExist).Once()
			},
		},
	}
	for _, t := range tests {
		if t.mock != nil {
			t.mock()
		}
		_, err := suite.gs.Load(t.ctx, (t.request).(*pb.LoadRequest))
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
