package postgres

import (
	"context"

	"github.com/kTowkA/GophKeeper/internal/model"
)

func (p *Postgres) Register(context.Context, model.StorageRegisterRequest) (model.StorageRegisterResponse, error) {
	return model.StorageRegisterResponse{}, nil
}
func (p *Postgres) PasswordHash(context.Context, model.StoragePasswordHashRequest) (model.StoragePasswordHashResponse, error) {
	return model.StoragePasswordHashResponse{}, nil
}
