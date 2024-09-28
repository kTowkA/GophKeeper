package models

import (
	pb "github.com/kTowkA/GophKeeper/grpc"
	"github.com/kTowkA/GophKeeper/internal/crypto"
)

// Services дополнительные сервисы, которые используеют наши модели. Они будут общими для всех
type Services struct {
	service pb.GophKeeperClient
	crypter crypto.Crypter
}

// NewServices создание переменной с сервисами
func NewServices(service pb.GophKeeperClient, crypter crypto.Crypter) *Services {
	return &Services{
		service: service,
		crypter: crypter,
	}
}
func (o *Services) SetRPCService(service pb.GophKeeperClient) {
	o.service = service
}
func (o *Services) RPCService() pb.GophKeeperClient {
	return o.service
}
func (o *Services) SetCrypter(crypter crypto.Crypter) {
	o.crypter = crypter
}
func (o *Services) Crypter() crypto.Crypter {
	return o.crypter
}
