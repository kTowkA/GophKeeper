package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kTowkA/GophKeeper/internal/model"
	"github.com/kTowkA/GophKeeper/internal/storage"
)

func (p *Postgres) Save(ctx context.Context, r model.StorageSaveRequest) (model.StorageSaveResponse, error) {
	userID, err := p.userID(ctx, r.User)
	if err != nil {
		return model.StorageSaveResponse{}, err
	}
	folderID, err := p.folderID(ctx, userID, r.Folder)
	if err != nil {
		return model.StorageSaveResponse{}, err
	}
	err = p.Pool.QueryRow(
		ctx,
		"SELECT value_id FROM keep_value WHERE folder_id=$1 AND title=$2",
		folderID,
		r.Value.Title,
	).Scan(nil)
	if err == nil {
		return model.StorageSaveResponse{}, storage.ErrKeepValueIsExist
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return model.StorageSaveResponse{}, err
	}
	valueID := uuid.New()
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return model.StorageSaveResponse{}, err
	}
	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO keep_value(value_id,folder_id,title,description,value,create_at) VALUES($1,$2,$3,$4,$5,$6);
		`,
		valueID,
		folderID,
		r.Value.Title,
		r.Value.Description,
		r.Value.Value,
		time.Now(),
	)
	if err != nil {
		return model.StorageSaveResponse{}, err
	}
	_, err = tx.Exec(
		ctx,
		`
		UPDATE keep_folder SET update_at=$2 WHERE folder_id=$1;
		`,
		folderID,
		time.Now(),
	)
	if err != nil {
		_ = tx.Rollback(ctx)
		return model.StorageSaveResponse{}, err
	}
	err = tx.Commit(ctx)
	return model.StorageSaveResponse{ValueID: valueID}, err
}
func (p *Postgres) Load(ctx context.Context, r model.StorageLoadRequest) (model.StorageLoadResponse, error) {
	userID, err := p.userID(ctx, r.User)
	if err != nil {
		return model.StorageLoadResponse{}, err
	}
	folderID, err := p.folderID(ctx, userID, r.Folder)
	if err != nil {
		return model.StorageLoadResponse{}, err
	}
	var value model.KeeperValue
	err = p.Pool.QueryRow(
		ctx,
		`
		SELECT 
			title,description,value,create_at
		FROM 
			keep_value
		WHERE 
			folder_id=$1 AND title=$2
		`,
		folderID,
		r.Title,
	).Scan(
		&value.Title,
		&value.Description,
		&value.Value,
		&value.CreateTime,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.StorageLoadResponse{}, storage.ErrKeepValueNotExist
	}
	if err != nil {
		return model.StorageLoadResponse{}, err
	}
	return model.StorageLoadResponse{Value: value}, nil
}
func (p *Postgres) CreateFolder(ctx context.Context, r model.StorageCreateFolderRequest) (model.StorageCreateFolderResponse, error) {
	userID, err := p.userID(ctx, r.User)
	if err != nil {
		return model.StorageCreateFolderResponse{}, err
	}
	err = p.Pool.QueryRow(
		ctx,
		"SELECT folder_id FROM keep_folder WHERE user_id=$1 AND title=$2",
		userID,
		r.Folder,
	).Scan(nil)
	if err == nil {
		return model.StorageCreateFolderResponse{}, storage.ErrKeepFolderIsExist
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return model.StorageCreateFolderResponse{}, err
	}
	folderID := uuid.New()
	_, err = p.Pool.Exec(
		ctx,
		"INSERT INTO keep_folder(folder_id,user_id,title,description,create_at,update_at) VALUES($1,$2,$3,$4,$5,$6) ",
		folderID,
		userID,
		r.Folder,
		r.Description,
		time.Now(),
		time.Now(),
	)
	return model.StorageCreateFolderResponse{FolderID: folderID}, err
}
func (p *Postgres) Folders(ctx context.Context, r model.StorageFoldersRequest) (model.StorageFoldersResponse, error) {
	userID, err := p.userID(ctx, r.User)
	if err != nil {
		return model.StorageFoldersResponse{}, err
	}
	rows, err := p.Pool.Query(
		ctx,
		"SELECT title FROM keep_folder WHERE user_id=$1",
		userID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.StorageFoldersResponse{}, storage.ErrKeepFolderNotExist
	}
	if err != nil {
		return model.StorageFoldersResponse{}, err
	}
	defer rows.Close()
	foldes := make([]string, 0, 1)
	for rows.Next() {
		folder := ""
		err = rows.Scan(&folder)
		if err != nil {
			return model.StorageFoldersResponse{}, err
		}
		foldes = append(foldes, folder)
	}
	if len(foldes) == 0 {
		return model.StorageFoldersResponse{}, storage.ErrKeepFolderNotExist
	}
	return model.StorageFoldersResponse{Folders: foldes}, nil
}
func (p *Postgres) Values(ctx context.Context, r model.StorageValuesRequest) (model.StorageValuesResponse, error) {
	userID, err := p.userID(ctx, r.User)
	if err != nil {
		return model.StorageValuesResponse{}, err
	}
	folderID, err := p.folderID(ctx, userID, r.Folder)
	if err != nil {
		return model.StorageValuesResponse{}, err
	}
	rows, err := p.Pool.Query(
		ctx,
		`
		SELECT 
			title
		FROM 
			keep_value
		WHERE 
			folder_id=$1
		`,
		folderID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.StorageValuesResponse{}, storage.ErrKeepValueNotExist
	}
	if err != nil {
		return model.StorageValuesResponse{}, err
	}
	defer rows.Close()
	values := make([]string, 0, 1)
	for rows.Next() {
		value := ""
		err = rows.Scan(
			&value,
		)
		if err != nil {
			return model.StorageValuesResponse{}, err
		}
		values = append(values, value)
	}
	if len(values) == 0 {
		return model.StorageValuesResponse{}, storage.ErrKeepValueNotExist
	}
	return model.StorageValuesResponse{Values: values}, nil
}
func (p *Postgres) folderID(ctx context.Context, userID uuid.UUID, folder string) (uuid.UUID, error) {
	var folderID uuid.UUID
	err := p.Pool.QueryRow(
		ctx,
		"SELECT folder_id FROM keep_folder WHERE user_id=$1 AND title=$2",
		userID,
		folder,
	).Scan(&folderID)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.UUID{}, storage.ErrKeepFolderNotExist
	}
	if err != nil {
		return uuid.UUID{}, err
	}
	return userID, nil
}
