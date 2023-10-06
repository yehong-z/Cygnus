package service

import (
	"github.com/yehong-z/Cygnus/like/service/api"

	"dubbo.apache.org/dubbo-go/v3/config"
)

type ThumbupServerImpl struct {
	api.UnimplementedThumbupServer
}

func init() {
	config.SetProviderService(&ThumbupServerImpl{})
}
