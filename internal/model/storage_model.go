package model

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
