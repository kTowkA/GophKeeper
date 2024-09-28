// В данном файле описаны интерфейсы хранилища. В идеальном мире можно было бы делать две базы данных отдельно для пользователей и самих данных
// все методы используют модели данных, таким образом при изменении модели не придется менять код в данной области
package storage

import (
	"context"

	"github.com/kTowkA/GophKeeper/internal/model"
)

// Storager храниилище общее
type Storager interface {
	UserStorager
	KeepStorager
}

// UserStorager для хранилища пользователей
type UserStorager interface {
	// Register регистрация пользователя
	Register(context.Context, model.StorageRegisterRequest) (model.StorageRegisterResponse, error)
	// PasswordHash получение хеша пароля пользователя для валидации
	PasswordHash(context.Context, model.StoragePasswordHashRequest) (model.StoragePasswordHashResponse, error)
}

// KeepStorager для хранилища данных
type KeepStorager interface {
	// CreateFolder создание директории у пользователя
	CreateFolder(context.Context, model.StorageCreateFolderRequest) (model.StorageCreateFolderResponse, error)
	// Folders получение списка с директориями пользователя
	Folders(context.Context, model.StorageFoldersRequest) (model.StorageFoldersResponse, error)
	// Values получение списка названий данных хранимых пользователем
	Values(context.Context, model.StorageValuesRequest) (model.StorageValuesResponse, error)
	// Save сохранение данных пользователя
	Save(context.Context, model.StorageSaveRequest) (model.StorageSaveResponse, error)
	// Load получение конкретных данных пользователя
	Load(context.Context, model.StorageLoadRequest) (model.StorageLoadResponse, error)
}
