package model

import (
	"time"

	"github.com/google/uuid"
)

type StorageRegisterRequest struct {
	Login    string
	Password string
}

type StorageRegisterResponse struct {
}

type StoragePasswordHashRequest struct {
	Login string
}
type StoragePasswordHashResponse struct {
	PasswordHash string
}
type KeeperFolder struct {
	Title       string
	Description string
	CreateTime  time.Time
	UpdateTime  time.Time
}

type KeeperValue struct {
	Title       string
	Description string
	Value       []byte
	CreateTime  time.Time
}
type StorageSaveRequest struct {
	User   string
	Folder string
	Value  KeeperValue
}
type StorageSaveResponse struct {
	ValueID uuid.UUID
}
type StorageCreateFolderRequest struct {
	User        string
	Folder      string
	Description string
}
type StorageCreateFolderResponse struct {
	FolderID uuid.UUID
}
type StorageFoldersRequest struct {
	User string
}
type StorageFoldersResponse struct {
	Folders []string
}
type StorageLoadRequest struct {
	User   string
	Folder string
	Title  string
}
type StorageLoadResponse struct {
	Value KeeperValue
}
type StorageValuesRequest struct {
	User   string
	Folder string
}
type StorageValuesResponse struct {
	Values []string
}
