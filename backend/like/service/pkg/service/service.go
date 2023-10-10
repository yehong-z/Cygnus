package service

import (
	"context"
	"fmt"
	"github.com/yehong-z/Cygnus/like/service/api"
	"github.com/yehong-z/Cygnus/like/service/pkg/dao"
	"sync"
)

type ThumbupServerImpl struct {
	api.UnimplementedThumbupServer
	dao dao.Dao
}

func NewThumbupServerImpl(dao dao.Dao) *ThumbupServerImpl {
	return &ThumbupServerImpl{dao: dao}
}

func (t *ThumbupServerImpl) Like(ctx context.Context, req *api.LikeReq) (*api.LikeResp, error) {
	// 先查询旧的点赞状态
	state, err := t.Stats(ctx, &api.StatsReq{UserId: req.UserId, ObjectId: req.ObjectId})
	if err != nil {
		return nil, err
	}

	newType := int64(req.Action)
	// 不存在则异步执行操作
	if state == nil {
		err = t.dao.AddLikeState(ctx, req.UserId, req.ObjectId, newType)
		if err != nil {
			return nil, err
		}
	} else {
		// 存在则对比状态异步处理
		err = t.dao.UpdateLikeState(ctx, req.UserId, req.ObjectId, newType)
		if err != nil {
			return nil, err
		}
	}

	// 更新点赞计数
	like, dislike := calculateCount(state, req.Action)
	if like != 0 || dislike != 0 {
		err = t.dao.AsyncAddCount(ctx, req.ObjectId, like, dislike)
		if err != nil {
			return nil, err
		}
	}

	return &api.LikeResp{
		UserId:       req.UserId,
		ObjectId:     req.ObjectId,
		LikeCount:    state.LikeCount + like,
		DislikeCount: state.DislikeCount + dislike,
	}, nil
}

func (t *ThumbupServerImpl) Stats(ctx context.Context, req *api.StatsReq) (*api.StatsResp, error) {
	state, err := t.dao.GetLikeState(ctx, req.UserId, req.ObjectId)
	if err != nil {
		return nil, err
	}

	count, err := t.dao.GetLikeCount(ctx, req.ObjectId)
	if err != nil {
		return nil, err
	}

	return &api.StatsResp{
		ObjectId:     state.ObjectID,
		UserId:       state.UserID,
		LikeCount:    count.LikesCount,
		DislikeCount: count.DislikesCount,
		LikeState:    api.State(state.Type),
	}, nil
}

func (t *ThumbupServerImpl) MultiStats(ctx context.Context, req *api.MultiStatsReq) (*api.MultiStatsResp, error) {
	stateMp, err := t.dao.GetLikeStateBatch(ctx, req.UserId, req.ObjectId)
	if err != nil {
		return nil, err
	}

	countMap, err := t.dao.GetLikeCountBatch(ctx, req.ObjectId)
	if err != nil {
		return nil, err
	}
	result := &api.MultiStatsResp{
		StatsResps: make([]*api.StatsResp, len(req.ObjectId)),
	}

	for i, obj := range req.ObjectId {
		result.StatsResps[i] = &api.StatsResp{
			ObjectId:     obj,
			UserId:       req.UserId,
			LikeState:    api.State(stateMp[obj].Type),
			LikeCount:    countMap[obj].LikesCount,
			DislikeCount: countMap[obj].DislikesCount,
		}
	}

	return result, nil
}

func (t *ThumbupServerImpl) HasLike(ctx context.Context, req *api.HasLikeReq) (*api.HasLikeResp, error) {
	state, err := t.dao.GetLikeState(ctx, req.UserId, req.ObjectId)
	if err != nil {
		return nil, err
	}

	return &api.HasLikeResp{
		UserLikeState: &api.UserLikeState{
			UserId: req.UserId,
			Time:   state.Mtime.Unix(),
			State:  api.State(state.Type),
		},
	}, nil
}

func (t *ThumbupServerImpl) UserLikes(ctx context.Context, req *api.UserLikesReq) (*api.UserLikesResp, error) {
	objs, err := t.dao.GetLikeStateBatch(ctx, req.UserId, req.ObjectIds)
	if err != nil {
		return nil, err
	}

	result := &api.UserLikesResp{
		UserLikeStates: make([]*api.UserLikeState, 0),
	}
	for _, v := range objs {
		result.UserLikeStates = append(result.UserLikeStates, &api.UserLikeState{
			UserId:   v.UserID,
			ObjectId: v.ObjectID,
			State:    api.State(v.Type),
			Time:     v.Mtime.Unix(),
		})
	}

	return result, nil
}

func (t *ThumbupServerImpl) ItemLikes(ctx context.Context, req *api.ItemLikesReq) (*api.ItemLikesResp, error) {
	users, err := t.dao.GetUsersByObject(ctx, req.ObjectId)
	if err != nil {
		return nil, err
	}

	return &api.ItemLikesResp{UserId: users}, nil
}

func (t *ThumbupServerImpl) UserLikedCounts(ctx context.Context, req *api.UserLikedCountsReq) (*api.UserLikedCountsResp, error) {
	cnt, err := t.dao.GetUserLikeCount(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &api.UserLikedCountsResp{Count: cnt}, nil
}

func (t *ThumbupServerImpl) BatchLikedCounts(ctx context.Context, req *api.BatchLikedCountsReq) (*api.BatchLikedCountsResp, error) {
	wg := sync.WaitGroup{}
	wg.Add(len(req.UserId))
	result := &api.BatchLikedCountsResp{Counts: make([]int64, len(req.UserId))}
	for i := range req.UserId {
		go func() {
			cnt, err := t.dao.GetUserLikeCount(ctx, req.UserId[i])
			if err != nil {
				fmt.Println(err.Error())
			}

			result.Counts[i] = cnt
			wg.Done()
		}()
	}

	wg.Wait()
	return result, nil
}

func calculateCount(stat *api.StatsResp, typ api.Action) (int64, int64) {
	var a, b, c, d int64 = 0, 0, 0, 0
	if stat.LikeState == api.State_STATE_LIKE {
		a = 1
	} else if stat.LikeState == api.State_STATE_DISLIKE {
		b = 1
	}

	if typ == api.Action_ACTION_LIKE {
		c = 1
	} else if typ == api.Action_ACTION_DISLIKE {
		d = 1
	}

	return c - a, d - b
}
