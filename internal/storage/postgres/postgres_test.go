// в данном файле содержится Suite-тест для всех реализованных методов Postgres
package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kTowkA/GophKeeper/internal/model"
	"github.com/kTowkA/GophKeeper/internal/storage"
	"github.com/kTowkA/GophKeeper/internal/storage/postgres/migrations"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type PostgresTest struct {
	suite.Suite
	pstorage *Postgres
	clear    dockerClear
}
type dockerClear struct {
	resource *dockertest.Resource
	pool     *dockertest.Pool
}

func (suite *PostgresTest) SetupSuite() {

	pool, err := dockertest.NewPool("")
	suite.Require().NoError(err)

	err = pool.Client.Ping()
	suite.Require().NoError(err)

	resource, err := pool.Run("postgres", "16", []string{"POSTGRES_USER=user", "POSTGRES_PASSWORD=pass"})
	suite.Require().NoError(err)

	connString := fmt.Sprintf("postgresql://user:pass@localhost:%s/user?sslmode=disable", resource.GetPort("5432/tcp"))

	// не нравится мне retry от dockertest, не всегда срабатывает нормально
	for range []int{1, 2, 3} {
		time.Sleep(5 * time.Second)
		err = func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			conn, err := pgx.Connect(ctx, connString)
			if err != nil {
				return err
			}
			defer conn.Close(ctx)
			return conn.Ping(ctx)
		}()
		if err == nil {
			break
		}
	}
	suite.Require().NoError(err)

	suite.clear = dockerClear{
		resource: resource,
		pool:     pool,
	}
	// ---------------------------------------------------------------------------------------------------

	suite.Require().NoError(err)
	err = migrations.MigrationsUP(connString)
	suite.Require().NoError(err)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ps, err := Connect(ctx, connString)
	suite.Require().NoError(err)
	suite.pstorage = ps

}
func (suite *PostgresTest) TearDownSuite() {
	suite.pstorage.Close()
	err := suite.clear.pool.Purge(suite.clear.resource)
	suite.NoError(err)
}
func (suite *PostgresTest) TestRegister() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := suite.pstorage.Register(ctx, model.StorageRegisterRequest{
		Login: "qwerty",
	})
	suite.NoError(err)
	_, err = suite.pstorage.Register(ctx, model.StorageRegisterRequest{
		Login: "QWERTY",
	})
	suite.ErrorIs(err, storage.ErrLoginIsAlreadyOccupied)
	_, err = suite.pstorage.Register(ctx, model.StorageRegisterRequest{
		Login: " ",
	})
	suite.Error(err)
}
func (suite *PostgresTest) TestPassword() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	login := "test_password"
	password := "test_password"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	suite.NoError(err)

	_, err = suite.pstorage.Register(ctx, model.StorageRegisterRequest{
		Login:    login,
		Password: string(hash),
	})
	suite.NoError(err)

	resp, err := suite.pstorage.PasswordHash(ctx, model.StoragePasswordHashRequest{
		Login: login,
	})
	suite.NoError(err)
	err = bcrypt.CompareHashAndPassword([]byte(resp.PasswordHash), []byte(password))
	suite.NoError(err)

	_, err = suite.pstorage.PasswordHash(ctx, model.StoragePasswordHashRequest{
		Login: login + "X",
	})
	suite.ErrorIs(err, storage.ErrUserIsNotExist)
}
func (suite *PostgresTest) TestUserID() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := "user_TestUserID"
	_, err := suite.pstorage.Register(ctx, model.StorageRegisterRequest{Login: user, Password: user})
	suite.NoError(err)

	_, err = suite.pstorage.userID(ctx, user)
	suite.NoError(err)

	_, err = suite.pstorage.userID(ctx, user+"x")
	suite.ErrorIs(err, storage.ErrUserIsNotExist)
}
func (suite *PostgresTest) TestCreateFolder() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := "user_TestCreateFolder"
	_, err := suite.pstorage.Register(ctx, model.StorageRegisterRequest{Login: user, Password: user})
	suite.NoError(err)

	folder := "folder_user_TestCreateFolder"

	// нет пользователя
	_, err = suite.pstorage.CreateFolder(ctx, model.StorageCreateFolderRequest{User: user + "x", Folder: folder})
	suite.ErrorIs(err, storage.ErrUserIsNotExist)

	// все хорошо
	_, err = suite.pstorage.CreateFolder(ctx, model.StorageCreateFolderRequest{User: user, Folder: folder})
	suite.NoError(err)

	// повтор
	_, err = suite.pstorage.CreateFolder(ctx, model.StorageCreateFolderRequest{User: user, Folder: folder})
	suite.ErrorIs(err, storage.ErrKeepFolderIsExist)
}
func (suite *PostgresTest) TestFolders() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := "user_TestFolders"
	_, err := suite.pstorage.Register(ctx, model.StorageRegisterRequest{Login: user, Password: user})
	suite.NoError(err)

	// нет пользователя
	_, err = suite.pstorage.Folders(ctx, model.StorageFoldersRequest{User: user + "x"})
	suite.ErrorIs(err, storage.ErrUserIsNotExist)

	//  нет папок
	_, err = suite.pstorage.Folders(ctx, model.StorageFoldersRequest{User: user})
	suite.ErrorIs(err, storage.ErrKeepFolderNotExist)

	// создаем 2 папки
	folders := []string{"folder_TestFolders_1", "folder_TestFolders_2"}
	_, err = suite.pstorage.CreateFolder(ctx, model.StorageCreateFolderRequest{User: user, Folder: folders[0]})
	suite.NoError(err)
	_, err = suite.pstorage.CreateFolder(ctx, model.StorageCreateFolderRequest{User: user, Folder: folders[1]})
	suite.NoError(err)

	// есть 2 папки
	resp, err := suite.pstorage.Folders(ctx, model.StorageFoldersRequest{User: user})
	suite.NoError(err)
	suite.EqualValues(folders, resp.Folders)
}

func (suite *PostgresTest) TestSave() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := "user_TestSave"
	_, err := suite.pstorage.Register(ctx, model.StorageRegisterRequest{Login: user, Password: user})
	suite.NoError(err)

	folder := "folder_TestSave"
	_, err = suite.pstorage.CreateFolder(ctx, model.StorageCreateFolderRequest{User: user, Folder: folder})
	suite.NoError(err)

	value := model.KeeperValue{
		Title:       "v1",
		Description: "d1",
		Value:       []byte("123"),
	}
	// нет пользователя
	_, err = suite.pstorage.Save(ctx, model.StorageSaveRequest{User: user + "x", Folder: folder, Value: value})
	suite.ErrorIs(err, storage.ErrUserIsNotExist)

	// нет папки
	_, err = suite.pstorage.Save(ctx, model.StorageSaveRequest{User: user, Folder: folder + "x", Value: value})
	suite.ErrorIs(err, storage.ErrKeepFolderNotExist)

	// все хорошо
	_, err = suite.pstorage.Save(ctx, model.StorageSaveRequest{User: user, Folder: folder, Value: value})
	suite.NoError(err)

	// повтор
	_, err = suite.pstorage.Save(ctx, model.StorageSaveRequest{User: user, Folder: folder, Value: value})
	suite.ErrorIs(err, storage.ErrKeepValueIsExist)
}
func (suite *PostgresTest) TestValues() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := "user_TestValues"
	_, err := suite.pstorage.Register(ctx, model.StorageRegisterRequest{Login: user, Password: user})
	suite.NoError(err)

	folder := "folder_TestValues"

	// нет пользователя
	_, err = suite.pstorage.Values(ctx, model.StorageValuesRequest{User: user + "x", Folder: folder})
	suite.ErrorIs(err, storage.ErrUserIsNotExist)

	//  нет папок
	_, err = suite.pstorage.Values(ctx, model.StorageValuesRequest{User: user, Folder: folder})
	suite.ErrorIs(err, storage.ErrKeepFolderNotExist)

	// создаем нашу папку
	_, err = suite.pstorage.CreateFolder(ctx, model.StorageCreateFolderRequest{User: user, Folder: folder})
	suite.NoError(err)

	// нет данных
	_, err = suite.pstorage.Values(ctx, model.StorageValuesRequest{User: user, Folder: folder})
	suite.ErrorIs(err, storage.ErrKeepValueNotExist)

	// сохраняем данные
	v := model.KeeperValue{
		Title: "value_TestValues",
	}
	_, err = suite.pstorage.Save(ctx, model.StorageSaveRequest{
		User:   user,
		Folder: folder,
		Value:  v,
	})
	suite.NoError(err)

	// есть данные
	resp, err := suite.pstorage.Values(ctx, model.StorageValuesRequest{User: user, Folder: folder})
	suite.NoError(err)
	suite.EqualValues([]string{v.Title}, resp.Values)
}

func (suite *PostgresTest) TestLoad() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := "user_TestLoad"
	_, err := suite.pstorage.Register(ctx, model.StorageRegisterRequest{Login: user, Password: user})
	suite.NoError(err)

	folder := "folder_TestLoad"

	valueTitle := "value_TestLoad"

	// нет пользователя
	_, err = suite.pstorage.Load(ctx, model.StorageLoadRequest{
		User:   user + "x",
		Folder: folder,
		Title:  valueTitle,
	})
	suite.ErrorIs(err, storage.ErrUserIsNotExist)

	//  нет папок
	_, err = suite.pstorage.Load(ctx, model.StorageLoadRequest{
		User:   user,
		Folder: folder,
		Title:  valueTitle,
	})
	suite.ErrorIs(err, storage.ErrKeepFolderNotExist)

	// создаем нашу папку
	_, err = suite.pstorage.CreateFolder(ctx, model.StorageCreateFolderRequest{User: user, Folder: folder})
	suite.NoError(err)

	// нет данных
	_, err = suite.pstorage.Load(ctx, model.StorageLoadRequest{
		User:   user,
		Folder: folder,
		Title:  valueTitle,
	})
	suite.ErrorIs(err, storage.ErrKeepValueNotExist)

	// сохраняем данные
	v := model.KeeperValue{
		Title: valueTitle,
	}
	_, err = suite.pstorage.Save(ctx, model.StorageSaveRequest{
		User:   user,
		Folder: folder,
		Value:  v,
	})
	suite.NoError(err)

	// есть данные
	resp, err := suite.pstorage.Load(ctx, model.StorageLoadRequest{
		User:   user,
		Folder: folder,
		Title:  valueTitle,
	})
	suite.NoError(err)
	v.CreateTime = resp.Value.CreateTime
	suite.EqualValues(v, resp.Value)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTest))
}
