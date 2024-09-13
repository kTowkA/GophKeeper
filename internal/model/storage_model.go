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
type KeeperElement struct {
	Title       string
	Description string
	Type        string
	Values      []KeeperValue
}
type KeeperValue struct {
	Title       string
	Description string
	Value       []byte
}
type StorageSaveRequest struct {
	User  string
	Value KeeperElement
}
type StorageSaveResponse struct {
}
type StorageLoadRequest struct {
	User               string
	TitleKeeperElement string
}
type StorageLoadResponse struct {
	TitleKeeperElement KeeperElement
}
