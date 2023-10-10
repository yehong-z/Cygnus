package dao

import (
	"context"
	"github.com/yehong-z/Cygnus/common/mq"
	"github.com/yehong-z/Cygnus/like/service/pkg/dao/gen/model/model"
	"github.com/yehong-z/Cygnus/like/service/pkg/dao/gen/model/query"
	"gorm.io/gorm"
	"time"
)

type Dao interface {
	GetLikeState(ctx context.Context, userId, objectId int64) (*model.Like, error)
	GetLikeStateBatch(ctx context.Context, userId int64, objectId []int64) (map[int64]*model.Like, error)
	AddLikeState(ctx context.Context, userId, ObjectId int64, action int64) error
	UpdateLikeState(ctx context.Context, userId, ObjectId int64, action int64) error
	GetLikeCount(ctx context.Context, ObjectId int64) (*model.Count, error)
	GetLikeCountBatch(ctx context.Context, ObjectId []int64) (map[int64]*model.Count, error)
	GetUserLikeCount(ctx context.Context, userId int64) (int64, error)
	GetObjectsByUser(ctx context.Context, userId int64) ([]int64, error)
	GetUsersByObject(ctx context.Context, objectId int64) ([]int64, error)
	AsyncAddCount(ctx context.Context, objectId, like, dislike int64) error
}

type dao struct {
	mysql    *gorm.DB
	producer mq.Producer
}

func NewDao(mysql *gorm.DB, producer mq.Producer) Dao {
	return &dao{mysql: mysql, producer: producer}
}

func (d *dao) UpdateLikeState(ctx context.Context, userId, ObjectId int64, action int64) error {
	q := query.Use(d.mysql)
	_, err := q.Like.WithContext(ctx).Where(q.Like.UserID.Eq(userId)).Where(q.Like.ObjectID.Eq(ObjectId)).Update(q.Like.Type, action)
	if err != nil {
		return err
	}

	return err
}

func (d *dao) AsyncAddCount(ctx context.Context, objectId, like, dislike int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *dao) GetObjectsByUser(ctx context.Context, userId int64) ([]int64, error) {
	q := query.Use(d.mysql)
	find, err := q.Like.WithContext(ctx).Where(q.Like.UserID.Eq(userId)).Limit(10).Find()
	if err != nil {
		return nil, err
	}

	result := make([]int64, len(find))
	for i, like := range find {
		result[i] = like.ObjectID
	}

	return result, nil
}

func (d *dao) GetUsersByObject(ctx context.Context, objectId int64) ([]int64, error) {
	q := query.Use(d.mysql)
	find, err := q.Like.WithContext(ctx).Where(q.Like.UserID.Eq(objectId)).Limit(10).Find()
	if err != nil {
		return nil, err
	}

	result := make([]int64, len(find))
	for i, like := range find {
		result[i] = like.UserID
	}

	return result, nil
}

func (d *dao) GetLikeCountBatch(ctx context.Context, ObjectId []int64) (map[int64]*model.Count, error) {
	q := query.Use(d.mysql)
	find, err := q.Count.WithContext(ctx).Where(q.Count.ObjectID.In(ObjectId...)).Find()
	if err != nil {
		return nil, err
	}

	result := make(map[int64]*model.Count)
	for _, like := range find {
		result[like.ObjectID] = like
	}

	return result, nil
}

func (d *dao) GetLikeStateBatch(ctx context.Context, userId int64, objectId []int64) (map[int64]*model.Like, error) {
	q := query.Use(d.mysql)
	find, err := q.Like.WithContext(ctx).Where(q.Like.UserID.Eq(userId)).Where(q.Like.ObjectID.In(objectId...)).Find()
	if err != nil {
		return nil, err
	}

	result := make(map[int64]*model.Like)
	for _, like := range find {
		result[like.ObjectID] = like
	}

	return result, nil
}

func (d *dao) GetUserLikeCount(ctx context.Context, userId int64) (int64, error) {
	q := query.Use(d.mysql)
	return q.Like.WithContext(ctx).Where(q.Like.UserID.Eq(userId)).Count()
}

func (d *dao) GetLikeCount(ctx context.Context, ObjectId int64) (*model.Count, error) {
	q := query.Use(d.mysql)
	return q.Count.WithContext(ctx).Where(q.Count.ObjectID.Eq(ObjectId)).First()
}

func (d *dao) AddLikeState(ctx context.Context, userId, ObjectId int64, action int64) error {
	like := model.Like{UserID: userId, ObjectID: ObjectId, Type: action, Ctime: time.Now(), Mtime: time.Now()}
	return query.Use(d.mysql).Like.WithContext(ctx).Create(&like)
}

func (d *dao) GetLikeState(ctx context.Context, userId, objectId int64) (*model.Like, error) {
	q := query.Use(d.mysql)
	return q.Like.WithContext(ctx).Where(q.Like.UserID.Eq(userId)).Where(q.Like.ObjectID.Eq(objectId)).First()
}
