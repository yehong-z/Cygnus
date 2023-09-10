package service

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/config"
	"github.com/yehong-z/Cygnus/reply/interface/api"
)

type ReplyServerImpl struct {
	api.UnimplementedReplyServer
}

func (r ReplyServerImpl) Add(ctx context.Context, req *api.AddReq) (*api.AddResp, error) {
	//TODO implement me
	panic("implement me")
}

func (r ReplyServerImpl) Like(ctx context.Context, req *api.LikeReq) (*api.CommonResp, error) {
	//TODO implement me
	panic("implement me")
}

func (r ReplyServerImpl) CancelLike(ctx context.Context, req *api.CancelLikeReq) (*api.CommonResp, error) {
	//TODO implement me
	panic("implement me")
}

func (r ReplyServerImpl) Dislike(ctx context.Context, req *api.DislikeReq) (*api.CommonResp, error) {
	//TODO implement me
	panic("implement me")
}

func (r ReplyServerImpl) CancelDislike(ctx context.Context, req *api.CancelDislikeReq) (*api.CommonResp, error) {
	//TODO implement me
	panic("implement me")
}

func (r ReplyServerImpl) Report(ctx context.Context, req *api.ReportReq) (*api.CommonResp, error) {
	//TODO implement me
	panic("implement me")
}

func (r ReplyServerImpl) MainList(ctx context.Context, req *api.MainListReq) (*api.ReplysResp, error) {
	//TODO implement me
	panic("implement me")
}

func (r ReplyServerImpl) Detail(ctx context.Context, req *api.DetailReq) (*api.ReplysResp, error) {
	//TODO implement me
	panic("implement me")
}

func init() {
	config.SetProviderService(&ReplyServerImpl{})
}
