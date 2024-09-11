package storage

import "errors"

var (
	ErrLoginIsAlreadyOccupied = errors.New("такой логин уже занят")
)
