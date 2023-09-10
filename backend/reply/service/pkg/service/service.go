package service

import (
	"context"
)

import (
	"github.com/yehong-z/Cygnus/reply/service/api"

	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
)

type GreeterServerImpl struct {
	api.UnimplementedGreeterServer
}

func (s *GreeterServerImpl) SayHello(ctx context.Context, in *api.HelloRequest) (*api.User, error) {
	logger.Infof("Dubbo-go GreeterProvider get user name = %s\n", in.Name)
	return &api.User{Name: "Hello " + in.Name, Id: "12345", Age: 21}, nil
}

func init() {
	config.SetProviderService(&GreeterServerImpl{})
}
