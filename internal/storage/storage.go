package storage

import (
	"context"

	"github.com/kTowkA/GophKeeper/internal/model"
)

type Storager interface {
	Register(ctx context.Context, r model.StorageRegisterRequest) (model.StorageRegisterResponse, error)
	PasswordHash(ctx context.Context, r model.StoragePasswordHashRequest) (model.StoragePasswordHashResponse, error)
}
