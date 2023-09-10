package service

import (
	"context"
	"github.com/yehong-z/Cygnus/reply/service/api"

	"dubbo.apache.org/dubbo-go/v3/config"
)

type GreeterServerImpl struct {
	api.UnimplementedGreeterServer
}

func (s *GreeterServerImpl) SayHello(ctx context.Context, in *api.HelloRequest) (*api.User, error) {
	return &api.User{Name: "Hello " + in.Name, Id: "12345", Age: 21}, nil
}

func init() {
	config.SetProviderService(&GreeterServerImpl{})
}
