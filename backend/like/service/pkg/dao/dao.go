package dao

import (
	"context"
	"github.com/yehong-z/Cygnus/like/service/pkg/dao/gen/model/model"
	"github.com/yehong-z/Cygnus/like/service/pkg/dao/gen/model/query"
	"gorm.io/gorm"
	"time"
)

type Dao interface {
	GetLikeState(ctx context.Context, userId, objectId int64) (*model.Like, error)
	AddLikeState(ctx context.Context, userId, ObjectId int64, action int64) error
	GetLikeCount(ctx context.Context, ObjectId int64) (*model.Count, error)
}

type dao struct {
	mysql *gorm.DB
}

func (d *dao) GetLikeCount(ctx context.Context, ObjectId int64) (*model.Count, error) {
	return query.Use(d.mysql).Count.WithContext(ctx).Where(query.Count.ObjectID.Eq(ObjectId)).First()
}

func (d *dao) AddLikeState(ctx context.Context, userId, ObjectId int64, action int64) error {
	like := model.Like{UserID: userId, ObjectID: ObjectId, Type: action, Ctime: time.Now(), Mtime: time.Now()}
	return query.Use(d.mysql).Like.WithContext(ctx).Create(&like)
}

func (d *dao) GetLikeState(ctx context.Context, userId, objectId int64) (*model.Like, error) {
	return query.Use(d.mysql).Like.WithContext(ctx).Where(query.Like.UserID.Eq(userId)).Where(query.Like.ObjectID.Eq(objectId)).First()
}

func New() (Dao, error) {
	m, err := NewMysql()
	if err != nil {
		return nil, err
	}

	return &dao{mysql: m}, nil
}
