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
	suite.ErrorIs(err, storage.ErrLoginIsNotExist)
}

func (suite *PostgresTest) TestSave() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// нет пользователя
	_, err := suite.pstorage.Save(ctx, model.StorageSaveRequest{User: "none"})
	suite.ErrorIs(err, storage.ErrLoginIsNotExist)

	user := "save_user"
	_, err = suite.pstorage.Register(ctx, model.StorageRegisterRequest{Login: user, Password: user})
	suite.NoError(err)

	// все хорошо
	_, err = suite.pstorage.Save(ctx,
		model.StorageSaveRequest{
			User: user,
			Value: model.KeeperElement{
				Title: "title1",
				Values: []model.KeeperValue{
					{
						Title:       "1",
						Description: "11",
						Value:       []byte("111"),
					},
					{
						Title:       "2",
						Description: "22",
						Value:       []byte("222"),
					},
				},
			},
		})
	suite.NoError(err)

	// повторяется пользователь и заголовок элемента
	_, err = suite.pstorage.Save(ctx,
		model.StorageSaveRequest{
			User: user,
			Value: model.KeeperElement{
				Title:  "title1",
				Values: []model.KeeperValue{},
			},
		})

	suite.ErrorIs(err, storage.ErrKeepElementIsExist)

	// повторяется элемент и заголовок значения
	_, err = suite.pstorage.Save(ctx,
		model.StorageSaveRequest{
			User: user,
			Value: model.KeeperElement{
				Title: "title2",
				Values: []model.KeeperValue{
					{
						Title:       "1",
						Description: "11",
						Value:       []byte("111"),
					},
					{
						Title:       "1",
						Description: "22",
						Value:       []byte("222"),
					},
				},
			},
		})
	suite.Error(err)
}
func TestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTest))
}
