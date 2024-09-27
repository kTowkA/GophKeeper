package options

import pb "github.com/kTowkA/GophKeeper/grpc"

type Options struct {
	service pb.GophKeeperClient
	cr      Crypter
}

func CreateOptions(service pb.GophKeeperClient, crypter Crypter) *Options {
	return &Options{
		service: service,
		cr:      crypter,
	}
}
func (o *Options) SetService(service pb.GophKeeperClient) {
	o.service = service
}
func (o *Options) Service() pb.GophKeeperClient {
	return o.service
}
func (o *Options) Crypter() Crypter {
	return o.cr
}
