package dao

import (
	"context"
	"github.com/yehong-z/Cygnus/common/snow"
	"github.com/yehong-z/Cygnus/like/service/pkg/dao/gen/model/model"
	"github.com/yehong-z/Cygnus/like/service/pkg/dao/gen/model/query"
	"gorm.io/gorm"
)

type Dao interface {
	GetCount(ctx context.Context, objectId int64) (*model.Count, error)
	UpdateCount(ctx context.Context, objectId, like, dislike int64) error
	CreateCount(ctx context.Context, objectId, like, dislike int64) error
}

type dao struct {
	mysql *gorm.DB
	snow  *snow.Snow
}

func NewDao(mysql *gorm.DB, snow *snow.Snow) Dao {
	return &dao{mysql: mysql, snow: snow}
}

func (d *dao) GetCount(ctx context.Context, objectId int64) (*model.Count, error) {
	q := query.Use(d.mysql)
	find, err := q.Count.WithContext(ctx).Where(q.Count.ObjectID.Eq(objectId)).Limit(1).Find()
	if err != nil {
		return nil, err
	}

	if len(find) > 0 {
		return find[0], nil
	}

	return nil, nil
}

func (d *dao) UpdateCount(ctx context.Context, objectId, like, dislike int64) error {
	q := query.Use(d.mysql)
	_, err := q.Count.WithContext(ctx).Where(q.Count.ObjectID.Eq(objectId)).
		UpdateSimple(q.Count.LikesCount.Add(like),
			q.Count.DislikesCount.Add(dislike))
	if err != nil {
		return err
	}

	return err
}

func (d *dao) CreateCount(ctx context.Context, objectId, like, dislike int64) error {
	q := query.Use(d.mysql)
	id, err := d.snow.GetSnowFlakeID()
	if err != nil {
		return err
	}
	return q.Count.WithContext(ctx).Create(&model.Count{
		ID:            id,
		ObjectID:      objectId,
		LikesCount:    like,
		DislikesCount: dislike,
	})
}
