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

func (p *Postgres) Save(ctx context.Context, r model.StorageSaveRequest) (model.StorageSaveResponse, error) {
	var userID uuid.UUID
	login := strings.ToLower(r.User)
	err := p.Pool.QueryRow(
		ctx,
		"SELECT user_id FROM users WHERE login=$1",
		login,
	).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.StorageSaveResponse{}, storage.ErrLoginIsNotExist
	}
	if err != nil {
		return model.StorageSaveResponse{}, err
	}
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return model.StorageSaveResponse{}, err
	}
	elementID, err := saveElement(ctx, tx, userID, r.Value)
	if err != nil {
		_ = tx.Rollback(ctx)
		return model.StorageSaveResponse{}, err
	}
	_, err = saveValues(ctx, tx, elementID, r.Value.Values)
	if err != nil {
		_ = tx.Rollback(ctx)
		return model.StorageSaveResponse{}, err
	}
	err = tx.Commit(ctx)
	return model.StorageSaveResponse{}, err
}
func (p *Postgres) Load(ctx context.Context, r model.StorageLoadRequest) (model.StorageLoadResponse, error) {
	login := strings.ToLower(r.User)
	resp := model.StorageLoadResponse{
		TitleKeeperElement: model.KeeperElement{
			Values: make([]model.KeeperValue, 0, 1),
		},
	}
	elementID := uuid.UUID{}
	err := p.Pool.QueryRow(
		ctx,
		`
		SELECT 
			ke.element_id, ke.title, ke.description
		FROM 
			users AS u,keep_element AS ke
		WHERE 
				u.login=$1
			AND
				u.user_id = ke.user_id
			AND
				ke.title=$2
		`,
		login,
		r.TitleKeeperElement,
	).Scan(
		&elementID,
		&resp.TitleKeeperElement.Title,
		&resp.TitleKeeperElement.Description,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.StorageLoadResponse{}, storage.ErrKeepElementNotExist
	}
	if err != nil {
		return model.StorageLoadResponse{}, err
	}
	rows, err := p.Pool.Query(
		ctx,
		"SELECT title,description,value FROM keep_value WHERE element_id=$1",
		elementID,
	)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return model.StorageLoadResponse{}, err
	}
	defer rows.Close()
	for rows.Next() {
		value := model.KeeperValue{}
		err = rows.Scan(
			&value.Title,
			&value.Description,
			&value.Value,
		)
		if err != nil {
			return model.StorageLoadResponse{}, err
		}
		resp.TitleKeeperElement.Values = append(resp.TitleKeeperElement.Values, value)
	}
	return resp, nil
}

func saveElement(ctx context.Context, tx pgx.Tx, userID uuid.UUID, element model.KeeperElement) (uuid.UUID, error) {
	err := tx.QueryRow(
		ctx,
		"SELECT element_id FROM keep_element WHERE user_id=$1 AND title=$2",
		userID,
		element.Title,
	).Scan(nil)
	if err == nil {
		return uuid.UUID{}, storage.ErrKeepElementIsExist
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return uuid.UUID{}, err
	}
	elementID := uuid.New()
	_, err = tx.Exec(
		ctx,
		"INSERT INTO keep_element (element_id,user_id,title,description,adding_at) VALUES($1,$2,$3,$4,$5) ",
		elementID,
		userID,
		element.Title,
		element.Description,
		time.Now(),
	)
	return elementID, err
}
func saveValues(ctx context.Context, tx pgx.Tx, elementID uuid.UUID, values []model.KeeperValue) ([]uuid.UUID, error) {
	b := pgx.Batch{}
	valuesIDs := make([]uuid.UUID, 0, len(values))
	for _, v := range values {
		valueID := uuid.New()
		valuesIDs = append(valuesIDs, valueID)
		b.Queue(
			"INSERT INTO keep_value(value_id,element_id,title,description,value,adding_at) VALUES($1,$2,$3,$4,$5,$6)",
			valueID,
			elementID,
			v.Title,
			v.Description,
			v.Value,
			time.Now(),
		)
	}
	br := tx.SendBatch(ctx, &b)
	err := br.Close()
	return valuesIDs, err
}
