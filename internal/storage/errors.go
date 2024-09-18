package storage

import "errors"

var (
	ErrLoginIsAlreadyOccupied = errors.New("такой логин уже занят")
	ErrLoginIsNotExist        = errors.New("пользователя с таким логином не существует")
	ErrKeepElementIsExist     = errors.New("такой элемент уже существует")
)
