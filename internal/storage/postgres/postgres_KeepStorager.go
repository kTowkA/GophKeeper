package postgres

import (
	"context"

	"github.com/kTowkA/GophKeeper/internal/model"
)

func (p *Postgres) Save(context.Context, model.StorageSaveRequest) (model.StorageSaveResponse, error) {
	return model.StorageSaveResponse{}, nil
}
func (p *Postgres) Load(context.Context, model.StorageLoadRequest) (model.StorageLoadResponse, error) {
	return model.StorageLoadResponse{}, nil
}
