package storage

import (
	"context"

	"github.com/kTowkA/GophKeeper/internal/model"
)

type Storager interface {
	UserStorager
	KeepStorager
}
type UserStorager interface {
	Register(context.Context, model.StorageRegisterRequest) (model.StorageRegisterResponse, error)
	PasswordHash(context.Context, model.StoragePasswordHashRequest) (model.StoragePasswordHashResponse, error)
}

type KeepStorager interface {
	Save(context.Context, model.StorageSaveRequest) (model.StorageSaveResponse, error)
	Load(context.Context, model.StorageLoadRequest) (model.StorageLoadResponse, error)
}
