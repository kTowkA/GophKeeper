package storage

import "errors"

var (
	ErrLoginIsAlreadyOccupied = errors.New("такой логин уже занят")
	ErrUserIsNotExist         = errors.New("пользователя с таким логином не существует")
	ErrKeepFolderIsExist      = errors.New("такая папка уже существует")
	ErrKeepFolderNotExist     = errors.New("такая папка не существует")
	ErrKeepValueIsExist       = errors.New("такое значение уже существует")
	ErrKeepValueNotExist      = errors.New("такое значение не существует")
)
