package postgres

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kTowkA/GophKeeper/internal/model"
	"github.com/kTowkA/GophKeeper/internal/storage"
)

func (p *Postgres) Register(ctx context.Context, r model.StorageRegisterRequest) (model.StorageRegisterResponse, error) {
	login := strings.TrimSpace(strings.ToLower(r.Login))
	if login == "" {
		return model.StorageRegisterResponse{}, errors.New("логин не может быть пустым")
	}
	err := p.Pool.QueryRow(
		ctx,
		"SELECT user_id FROM users WHERE login=$1",
		login,
	).Scan(nil)
	if err == nil {
		return model.StorageRegisterResponse{}, storage.ErrLoginIsAlreadyOccupied
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return model.StorageRegisterResponse{}, err
	}
	_, err = p.Pool.Exec(
		ctx,
		"INSERT INTO users(user_id,login,password_hash,adding_at) VALUES($1,$2,$3,$4)",
		uuid.New(),
		login,
		r.Password,
		time.Now(),
	)
	if err != nil {
		return model.StorageRegisterResponse{}, err
	}
	return model.StorageRegisterResponse{}, nil
}
func (p *Postgres) PasswordHash(ctx context.Context, r model.StoragePasswordHashRequest) (model.StoragePasswordHashResponse, error) {
	login := strings.ToLower(r.Login)
	passwordHash := ""
	err := p.Pool.QueryRow(
		ctx,
		"SELECT password_hash FROM users WHERE login=$1",
		login,
	).Scan(&passwordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.StoragePasswordHashResponse{}, storage.ErrLoginIsNotExist
	}
	if err != nil {
		return model.StoragePasswordHashResponse{}, err
	}
	return model.StoragePasswordHashResponse{
		PasswordHash: passwordHash,
	}, nil
}
