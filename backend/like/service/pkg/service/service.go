package service

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/config"
	"github.com/yehong-z/Cygnus/like/service/api"
	"github.com/yehong-z/Cygnus/like/service/pkg/dao"
	"github.com/yehong-z/Cygnus/like/service/pkg/dao/gen/model/model"
)

type ThumbupServerImpl struct {
	api.UnimplementedThumbupServer
	dao dao.Dao
}

func (t *ThumbupServerImpl) Like(ctx context.Context, req *api.LikeReq) (*api.LikeResp, error) {
	// 先查询旧的点赞状态
	state, err := t.dao.GetLikeState(ctx, req.UserId, req.ObjectId)
	if err != nil {
		return nil, err
	}

	count, err := t.dao.GetLikeCount(ctx, req.ObjectId)
	newType := int64(req.Action)
	// 不存在则异步执行操作
	if state == nil {
		err = t.dao.AddLikeState(ctx, req.UserId, req.ObjectId, newType)
		if err != nil {
			return nil, err
		}
	} else {
		// 存在则对比状态异步处理

	}

	// 更新点赞计数
	cnt := calculateCount(state, count, req.Action)
	return &api.LikeResp{UserId: req.UserId, ObjectId: req.ObjectId, LikeCount: cnt.LikesCount, DislikeCount: cnt.DislikesCount}, nil
}

func (t *ThumbupServerImpl) Stats(ctx context.Context, req *api.StatsReq) (*api.StatsResp, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThumbupServerImpl) MultiStats(ctx context.Context, req *api.MultiStatsReq) (*api.MultiStatsResp, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThumbupServerImpl) HasLike(ctx context.Context, req *api.HasLikeReq) (*api.HasLikeResp, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThumbupServerImpl) UserLikes(ctx context.Context, req *api.UserLikesReq) (*api.UserLikesResp, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThumbupServerImpl) ItemLikes(ctx context.Context, req *api.ItemLikesReq) (*api.ItemLikesResp, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThumbupServerImpl) UserLikedCounts(ctx context.Context, req *api.UserLikedCountsReq) (*api.UserLikedCountsResp, error) {
	//TODO implement me
	panic("implement me")
}

func (t *ThumbupServerImpl) BatchLikedCounts(ctx context.Context, req *api.BatchLikedCountsReq) (*api.BatchLikedCountsResp, error) {
	//TODO implement me
	panic("implement me")
}

func calculateCount(stat *model.Like, count *model.Count, typ api.Action) *model.Count {
	var likesCount, dislikesCount int64
	switch typ {
	case api.Action_ACTION_LIKE:
		likesCount = 1
		if stat.Type == 2 {
			dislikesCount = -1
		}
	case api.Action_ACTION_CANCEL_LIKE:
		likesCount = -1
	case api.Action_ACTION_DISLIKE:
		dislikesCount = 1
		if stat.Type == 1 {
			likesCount = -1
		}
	case api.Action_ACTION_CANCEL_DISLIKE:
		dislikesCount = -1
	}
	count.LikesCount += likesCount
	count.DislikesCount += dislikesCount
	return count
}

func init() {
	config.SetProviderService(&ThumbupServerImpl{})
}
