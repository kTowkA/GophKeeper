// содержит модели данных для общения сервера с хранилищем и данных реализующих сообщения в proto-файле
package model

import (
	"time"

	"github.com/google/uuid"
)

// StorageRegisterRequest запрос в хранилище регистраици пользователя
type StorageRegisterRequest struct {
	// Login логин пользователя
	Login string
	// Password пароль пользователя
	Password string
}

// StorageRegisterResponse ответ при запросе в хранилище регистрации пользователя
type StorageRegisterResponse struct {
}

// StoragePasswordHashRequest запрос в хранилище хеша пароля пользователя
type StoragePasswordHashRequest struct {
	// Login логин пользователя
	Login string
}

// StoragePasswordHashResponse ответ при запросе в хранилище хеша пароля
type StoragePasswordHashResponse struct {
	// PasswordHash хеш пароля
	PasswordHash string
}

// KeeperFolder описание директории (как в proto)
type KeeperFolder struct {
	// Title название директории
	Title string
	// Description описание директории
	Description string
	// CreateTime временная метка создания
	CreateTime time.Time
	// UpdateTime временная метка обновления
	UpdateTime time.Time
}

// KeeperValue описание структуры с сохраняемыми данными (как в proto)
type KeeperValue struct {
	// Title название данных
	Title string
	// Description описание данных
	Description string
	// Value собственно значение
	Value []byte
	// CreateTime временная метка создания
	CreateTime time.Time
}

// StorageSaveRequest запрос в хранилище сохранения данных
type StorageSaveRequest struct {
	// User логин пользователя
	User string
	// Folder директория
	Folder string
	// Value значение
	Value KeeperValue
}

// StorageSaveResponse ответ при запросе в хранилище сохранения данных
type StorageSaveResponse struct {
	// ValueID UUID сохраненного значения
	ValueID uuid.UUID
}

// StorageCreateFolderRequest запрос в хранилище создание директории
type StorageCreateFolderRequest struct {
	// User логин пользователя
	User string
	// Folder директория
	Folder string
	// Description описание директории
	Description string
}

// StorageCreateFolderResponse ответ при запросе в хранилище создание директории
type StorageCreateFolderResponse struct {
	// FolderID UUID созданной директории
	FolderID uuid.UUID
}

// StorageFoldersRequest запрос в хранилище списка директорий
type StorageFoldersRequest struct {
	// User логин пользователя
	User string
}

// StorageFoldersResponse ответ при запросе в хранилище списка директорий
type StorageFoldersResponse struct {
	// Folders список наименований директорий
	Folders []string
}

// StorageLoadRequest запрос в хранилище получения данных
type StorageLoadRequest struct {
	// User логин пользователя
	User string
	// Folder директория
	Folder string
	// Title название данных
	Title string
}

// StorageLoadResponse ответ при запросе в хранилище данных
type StorageLoadResponse struct {
	// Value данные
	Value KeeperValue
}

// StorageValuesRequest запрос в хранилище сохранения данных
type StorageValuesRequest struct {
	// User логин пользователя
	User string
	// Folder директория
	Folder string
}

// StorageValuesResponse ответ при запросе в хранилище списка данных
type StorageValuesResponse struct {
	// Values список наименований данных
	Values []string
}
