package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/yehong-z/Cygnus/like/service/api"
	"github.com/yehong-z/Cygnus/like/service/pkg/dao"
	"github.com/yehong-z/Cygnus/like/service/pkg/dao/gen/model/model"
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
	if err != nil || state == nil {
		return nil, errors.New("get state error")
	}

	newSt, like, dislike := calculateCount(state, req.Action)
	// 不存在则异步执行操作
	if state.UserId != req.UserId || state.ObjectId != req.ObjectId {
		err = t.dao.AddLikeState(ctx, req.UserId, req.ObjectId, int64(newSt.LikeState))
	} else {
		err = t.dao.UpdateLikeState(ctx, req.UserId, req.ObjectId, int64(newSt.LikeState))
	}

	if err != nil {
		return nil, err
	}

	// 更新点赞计数
	if like != 0 || dislike != 0 {
		err = t.dao.AsyncAddCount(ctx, req.ObjectId, like, dislike)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return &api.LikeResp{
		UserId:       req.UserId,
		ObjectId:     req.ObjectId,
		LikeCount:    newSt.LikeCount,
		DislikeCount: newSt.DislikeCount,
	}, nil
}

func (t *ThumbupServerImpl) Stats(ctx context.Context, req *api.StatsReq) (*api.StatsResp, error) {
	state, err := t.dao.GetLikeState(ctx, req.UserId, req.ObjectId)
	if err != nil || state == nil {
		state = &model.Like{}
	}

	count, err := t.dao.GetLikeCount(ctx, req.ObjectId)
	if err != nil || count == nil {
		count = &model.Count{}
	}

	result := &api.StatsResp{
		ObjectId:     state.ObjectID,
		UserId:       state.UserID,
		LikeState:    api.State(state.Type),
		LikeCount:    count.LikesCount,
		DislikeCount: count.DislikesCount,
	}

	return result, nil
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

func calculateCount(stat *api.StatsResp, typ api.Action) (*api.StatsResp, int64, int64) {
	var a, b, c, d int64 = 0, 0, 0, 0
	if stat.LikeState == api.State_STATE_LIKE {
		a = 1
	} else if stat.LikeState == api.State_STATE_DISLIKE {
		b = 1
	}

	likeSt := api.State_STATE_UNSPECIFIED
	if typ == api.Action_ACTION_LIKE {
		c = 1
		likeSt = api.State_STATE_LIKE
	} else if typ == api.Action_ACTION_DISLIKE {
		d = 1
		likeSt = api.State_STATE_DISLIKE
	} else if typ == api.Action_ACTION_CANCEL_DISLIKE && stat.LikeState != api.State_STATE_DISLIKE {
		return stat, 0, 0
	} else if typ == api.Action_ACTION_CANCEL_LIKE && stat.LikeState != api.State_STATE_LIKE {
		return stat, 0, 0
	}

	st := &api.StatsResp{
		ObjectId:     stat.ObjectId,
		UserId:       stat.UserId,
		LikeCount:    stat.LikeCount + (c - a),
		DislikeCount: stat.DislikeCount + (d - b),
		LikeState:    likeSt,
	}

	fmt.Println(a, b, c, d)
	return st, c - a, d - b
}
